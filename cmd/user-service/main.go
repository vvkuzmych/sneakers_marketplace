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

	"github.com/vvkuzmych/sneakers_marketplace/internal/user/handler"
	"github.com/vvkuzmych/sneakers_marketplace/internal/user/repository"
	"github.com/vvkuzmych/sneakers_marketplace/internal/user/service"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/auth"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/config"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/database"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/logger"
	pb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/user"
)

func main() {
	// Initialize logger
	log := logger.New(logger.Config{
		Level:  "debug",
		Format: "console",
		Output: os.Stdout,
	})
	log.Info("Starting User Service...")

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

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.Expiration,
		cfg.JWT.RefreshExpiration,
	)
	log.Info("JWT manager initialized")

	// Initialize repository
	userRepo := repository.NewUserRepository(db)

	// Initialize service
	userService := service.NewUserService(userRepo, jwtManager)

	// Initialize gRPC handler
	userHandler := handler.NewUserHandler(userService)

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor(log)),
	)

	// Register service
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	// Enable gRPC reflection for debugging
	reflection.Register(grpcServer)

	// Start gRPC server
	address := fmt.Sprintf(":%d", cfg.Server.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", address, err)
	}

	log.Infof("User Service listening on %s", address)

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

	log.Info("Shutting down User Service...")
	grpcServer.GracefulStop()
	log.Info("User Service stopped")
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
