package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/vvkuzmych/sneakers_marketplace/internal/gateway/clients"
	"github.com/vvkuzmych/sneakers_marketplace/internal/gateway/handlers"
	"github.com/vvkuzmych/sneakers_marketplace/internal/gateway/middleware"
	"github.com/vvkuzmych/sneakers_marketplace/internal/gateway/websocket"
)

// SetupRouter configures all routes
func SetupRouter(grpcClients *clients.GRPCClients, wsHub *websocket.Hub) *gin.Engine {
	router := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // TODO: Configure for production
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":         "healthy",
			"service":        "api-gateway",
			"ws_connections": wsHub.GetClientCount(),
		})
	})

	// WebSocket endpoint (requires JWT authentication via query param or header)
	router.GET("/ws", websocket.HandleWebSocket(wsHub))

	// Initialize handlers
	userHandler := handlers.NewUserHandler(grpcClients.UserClient)
	productHandler := handlers.NewProductHandler(grpcClients.ProductClient)
	biddingHandler := handlers.NewBiddingHandler(grpcClients.BiddingClient)
	orderHandler := handlers.NewOrderHandler(grpcClients.OrderClient)
	paymentHandler := handlers.NewPaymentHandler(grpcClients.PaymentClient)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// User routes (protected)
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/:user_id", userHandler.GetProfile)
		}

		// Product routes (public for read, protected for write)
		products := v1.Group("/products")
		{
			products.GET("", productHandler.ListProducts)
			products.GET("/:id", productHandler.GetProduct)
			products.GET("/search", productHandler.SearchProducts)
		}

		// Bidding routes (protected)
		bidding := v1.Group("")
		bidding.Use(middleware.AuthMiddleware())
		{
			bidding.POST("/bids", biddingHandler.PlaceBid)
			bidding.POST("/asks", biddingHandler.PlaceAsk)
			bidding.GET("/bids/product/:product_id", biddingHandler.GetProductBids)
			bidding.GET("/asks/product/:product_id", biddingHandler.GetProductAsks)
		}

		// Market routes (public)
		market := v1.Group("/market")
		{
			market.GET("/:product_id/:size_id", biddingHandler.GetMarketPrice)
		}

		// Order routes (protected)
		orders := v1.Group("/orders")
		orders.Use(middleware.AuthMiddleware())
		{
			orders.GET("/:id", orderHandler.GetOrder)
			orders.GET("/buyer/:buyer_id", orderHandler.ListBuyerOrders)
		}

		// Payment routes (protected)
		payments := v1.Group("/payments")
		payments.Use(middleware.AuthMiddleware())
		{
			payments.POST("/intent", paymentHandler.CreatePaymentIntent)
			payments.GET("/:id", paymentHandler.GetPayment)
		}
	}

	return router
}
