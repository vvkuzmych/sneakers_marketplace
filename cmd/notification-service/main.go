package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/vvkuzmych/sneakers_marketplace/internal/notification/email"
	"github.com/vvkuzmych/sneakers_marketplace/internal/notification/handler"
	"github.com/vvkuzmych/sneakers_marketplace/internal/notification/repository"
	"github.com/vvkuzmych/sneakers_marketplace/internal/notification/service"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/config"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/database"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/logger"
	pb "github.com/vvkuzmych/sneakers_marketplace/pkg/proto/notification"
)

func main() {
	// Initialize logger
	log := logger.NewDevelopment()

	log.Info("Starting Notification Service")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Initialize database
	ctx := context.Background()
	db, err := database.NewPostgresPool(ctx, database.PostgresConfig{
		URL: cfg.Database.URL,
	}, log)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	log.Info("Connected to PostgreSQL")

	// Initialize services
	notifRepo := repository.NewNotificationRepository(db)
	emailService := email.NewEmailService()
	notifService := service.NewNotificationService(notifRepo, emailService)
	notifHandler := handler.NewNotificationHandler(notifService)

	// Get port
	port := os.Getenv("NOTIFICATION_SERVICE_PORT")
	if port == "" {
		port = "50056"
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, notifHandler)

	// Enable gRPC reflection for debugging
	reflection.Register(grpcServer)

	// Start server
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err.Error())
	}

	go func() {
		log.Info("Notification Service listening on port " + port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err.Error())
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Notification Service")
	grpcServer.GracefulStop()
	log.Info("Notification Service stopped")
}
