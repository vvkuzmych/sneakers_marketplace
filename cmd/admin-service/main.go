package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/vvkuzmych/sneakers_marketplace/internal/admin/handler"
	"github.com/vvkuzmych/sneakers_marketplace/internal/admin/repository"
	"github.com/vvkuzmych/sneakers_marketplace/internal/admin/service"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/config"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/database"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/logger"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/middleware"
	adminpb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/admin"
)

func main() {
	// Initialize logger
	log := logger.New(logger.Config{
		Level:  "debug",
		Format: "console",
		Output: os.Stdout,
	})
	log.Info("🚀 Starting Admin Service...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Get JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET environment variable is required")
	}

	// Get port from environment variable or use default
	port := os.Getenv("ADMIN_SERVICE_PORT")
	if port == "" {
		port = "50057"
	}

	// Connect to database
	log.Info("📦 Connecting to PostgreSQL...")
	ctx := context.Background()
	db, err := database.NewPostgresPool(ctx, database.PostgresConfig{
		URL: cfg.Database.URL,
	}, log)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Info("✅ Connected to PostgreSQL")

	// Initialize repository, service, and handler
	adminRepo := repository.NewAdminRepository(db)
	adminService := service.NewAdminService(adminRepo)
	adminHandler := handler.NewAdminHandler(adminService)

	// Create gRPC server with RBAC middleware
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.LoggingInterceptor,
			middleware.RequireAdmin(jwtSecret),
		),
	)

	// Register services
	adminpb.RegisterAdminServiceServer(grpcServer, adminHandler)

	// Enable reflection for grpcurl
	reflection.Register(grpcServer)

	// Start listening
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	log.Infof("✅ Admin Service listening on :%s", port)
	log.Info("🔒 RBAC: Admin role required for all endpoints")
	log.Info("📝 gRPC Reflection enabled")

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Info("🛑 Shutting down Admin Service...")
		grpcServer.GracefulStop()
	}()

	// Start server
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
