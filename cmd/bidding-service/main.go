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

	"github.com/vvkuzmych/sneakers_marketplace/internal/bidding/handler"
	"github.com/vvkuzmych/sneakers_marketplace/internal/bidding/repository"
	"github.com/vvkuzmych/sneakers_marketplace/internal/bidding/service"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/config"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/database"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/logger"
	pb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/bidding"
)

func main() {
	// Initialize logger
	log := logger.New(logger.Config{
		Level:  "debug",
		Format: "console",
		Output: os.Stdout,
	})
	log.Info("Starting Bidding Service...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
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

	// Initialize repository
	biddingRepo := repository.NewBiddingRepository(db)

	// Initialize service
	biddingService := service.NewBiddingService(biddingRepo)

	// Initialize gRPC handler
	biddingHandler := handler.NewBiddingHandler(biddingService)

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor(log)),
	)

	// Register service
	pb.RegisterBiddingServiceServer(grpcServer, biddingHandler)

	// Enable gRPC reflection for debugging
	reflection.Register(grpcServer)

	// Get port from environment (Bidding Service uses 50053)
	// Use BIDDING_SERVICE_PORT env var, or default to 50053
	port := os.Getenv("BIDDING_SERVICE_PORT")
	if port == "" {
		port = "50053"
	}
	address := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", address, err)
	}

	log.Infof("Bidding Service listening on %s", address)

	// Graceful shutdown
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Bidding Service...")
	grpcServer.GracefulStop()
	log.Info("Bidding Service stopped")
}

// loggingInterceptor logs all gRPC requests
func loggingInterceptor(log *logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Call the handler
		resp, err := handler(ctx, req)

		// Log request
		duration := time.Since(start)
		if err != nil {
			log.Errorf("gRPC request failed: method=%s duration=%v error=%v", info.FullMethod, duration, err)
		} else {
			log.Infof("gRPC request completed: method=%s duration=%v", info.FullMethod, duration)
		}

		return resp, err
	}
}
