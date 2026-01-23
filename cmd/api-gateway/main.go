package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vvkuzmych/sneakers_marketplace/internal/gateway/clients"
	"github.com/vvkuzmych/sneakers_marketplace/internal/gateway/router"
	"github.com/vvkuzmych/sneakers_marketplace/internal/gateway/websocket"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/database"
	"github.com/vvkuzmych/sneakers_marketplace/pkg/logger"
)

func main() {
	log.Println("🚀 Starting API Gateway...")

	// Initialize logger
	logger := logger.New(logger.Config{
		Level:  "debug",
		Format: "console",
		Output: os.Stdout,
	})

	// Initialize database connection for fee tracking
	ctx := context.Background()
	databaseURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/sneakers_marketplace?sslmode=disable")
	db, err := database.NewPostgresPool(ctx, database.PostgresConfig{
		URL: databaseURL,
	}, logger)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("✅ Database connected")

	// Get service addresses from environment
	userAddr := getEnv("USER_SERVICE_ADDR", "localhost:50051")
	productAddr := getEnv("PRODUCT_SERVICE_ADDR", "localhost:50052")
	biddingAddr := getEnv("BIDDING_SERVICE_ADDR", "localhost:50053")
	orderAddr := getEnv("ORDER_SERVICE_ADDR", "localhost:50054")
	paymentAddr := getEnv("PAYMENT_SERVICE_ADDR", "localhost:50055")

	// Connect to all gRPC services
	grpcClients, err := clients.NewGRPCClients(
		userAddr, productAddr, biddingAddr, orderAddr, paymentAddr,
	)
	if err != nil {
		log.Fatalf("Failed to connect to gRPC services: %v", err)
	}

	// Create WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()
	log.Println("✅ WebSocket Hub started")

	// Setup router
	r := router.SetupRouter(grpcClients, wsHub, db, logger)

	// Get HTTP port
	port := getEnv("HTTP_PORT", "8080")

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("✅ API Gateway listening on :%s", port)
		log.Printf("📝 Health check: http://localhost:%s/health", port)
		log.Printf("📚 API endpoints: http://localhost:%s/api/v1", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down API Gateway...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ API Gateway stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
