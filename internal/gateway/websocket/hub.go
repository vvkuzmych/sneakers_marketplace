package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to clients
type Hub struct {
	// Registered clients (userID -> Client)
	clients map[int64]*Client

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[int64]*Client),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			// Disconnect existing client if same user
			if existingClient, ok := h.clients[client.UserID]; ok {
				existingClient.SafeClose()
				delete(h.clients, client.UserID)
			}
			h.clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("Client connected: UserID=%d, Total=%d", client.UserID, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				client.SafeClose()
				log.Printf("Client disconnected: UserID=%d, Total=%d", client.UserID, len(h.clients))
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			// Broadcast to all clients
			h.mu.RLock()
			for userID, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					// Client's send buffer is full, disconnect
					client.SafeClose()
					delete(h.clients, userID)
					log.Printf("Client send buffer full, disconnecting: UserID=%d", userID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// SendToUser sends a message to a specific user
func (h *Hub) SendToUser(userID int64, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()

	if !ok {
		// User not connected, skip
		return nil
	}

	select {
	case client.Send <- data:
		log.Printf("Sent message to UserID=%d", userID)
	default:
		log.Printf("Failed to send to UserID=%d (buffer full)", userID)
	}

	return nil
}

// BroadcastToAll sends a message to all connected clients
func (h *Hub) BroadcastToAll(message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	h.broadcast <- data
	return nil
}

// GetConnectedUserIDs returns list of connected user IDs
func (h *Hub) GetConnectedUserIDs() []int64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	userIDs := make([]int64, 0, len(h.clients))
	for userID := range h.clients {
		userIDs = append(userIDs, userID)
	}
	return userIDs
}

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// IsUserConnected checks if a user is connected
func (h *Hub) IsUserConnected(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}
