package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/vvkuzmych/sneakers_marketplace/internal/payment/handler"
	"github.com/vvkuzmych/sneakers_marketplace/internal/payment/repository"
	"github.com/vvkuzmych/sneakers_marketplace/internal/payment/service"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/config"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/database"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/logger"
	pb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/payment"
)

// loggingInterceptor logs all gRPC requests
func loggingInterceptor(log *logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)

		if err != nil {
			log.Errorf("gRPC request failed: method=%s duration=%v error=%v", info.FullMethod, duration, err)
		} else {
			log.Infof("gRPC request completed: method=%s duration=%v", info.FullMethod, duration)
		}

		return resp, err
	}
}

func main() {
	// Initialize logger
	log := logger.New(logger.Config{
		Level:  "debug",
		Format: "console",
		Output: os.Stdout,
	})
	log.Info("Starting Payment Service...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Override port for payment service
	port := os.Getenv("PAYMENT_SERVICE_PORT")
	if port != "" {
		portInt := 50055
		fmt.Sscanf(port, "%d", &portInt)
		cfg.Server.Port = portInt
	} else {
		cfg.Server.Port = 50055 // Default port for Payment Service
	}

	// Initialize database
	ctx := context.Background()
	db, err := database.NewPostgresPool(ctx, database.PostgresConfig{
		URL: cfg.Database.URL,
	}, log)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Info("Connected to database")

	// Initialize repository, service, and handler
	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor(log)),
	)

	// Register service
	pb.RegisterPaymentServiceServer(grpcServer, paymentHandler)

	// Enable gRPC reflection for debugging
	reflection.Register(grpcServer)

	// Start gRPC server
	address := fmt.Sprintf(":%d", cfg.Server.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", address, err)
	}

	log.Infof("Payment Service is running on %s", address)

	// Start serving in a goroutine
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Payment Service...")
	grpcServer.GracefulStop()
	log.Info("Payment Service stopped")
}
