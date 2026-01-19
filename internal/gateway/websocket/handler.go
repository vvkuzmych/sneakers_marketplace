package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Configure allowed origins for production
		return true // Allow all origins for development
	},
}

// HandleWebSocket handles WebSocket connections with JWT authentication
func HandleWebSocket(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get JWT token from query parameter or header
		token := c.Query("token")
		if token == "" {
			token = c.GetHeader("Authorization")
			if len(token) > 7 && token[:7] == "Bearer " {
				token = token[7:]
			}
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
			return
		}

		// Validate JWT token
		userID, email, err := validateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade connection: %v", err)
			return
		}

		// Create client
		client := &Client{
			UserID: userID,
			Email:  email,
			Hub:    hub,
			Conn:   conn,
			Send:   make(chan []byte, 256),
		}

		// Register client
		hub.register <- client

		// Send welcome message
		welcomeMsg := Message{
			Type: "connected",
			Data: map[string]interface{}{
				"user_id": userID,
				"email":   email,
				"message": "Connected to notifications",
			},
		}
		if data, err := json.Marshal(welcomeMsg); err == nil {
			client.Send <- data
		}

		// Start client goroutines
		go client.writePump()
		go client.readPump()
	}
}

// validateToken validates JWT token and extracts user info
func validateToken(tokenString string) (int64, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return 0, "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", fmt.Errorf("invalid claims")
	}

	// Extract user_id
	var userID int64
	switch v := claims["user_id"].(type) {
	case float64:
		userID = int64(v)
	case string:
		userID, _ = strconv.ParseInt(v, 10, 64)
	default:
		return 0, "", fmt.Errorf("invalid user_id type")
	}

	email, _ := claims["email"].(string)

	return userID, email, nil
}
