package hub

import (
	"butter-socket/models"
	"context"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a connected customer
type Client struct {
	Hub          *Hub
	Conn         *websocket.Conn
	Send         chan []byte
	Customer     *models.Customer
	Conversation *models.Conversation
	CancelAI     context.CancelFunc
}

// Hub maintains active clients and broadcasts messages
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from clients
	broadcast chan []byte

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// NewHub creates a new Hub instanceinstance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			fmt.Printf("Client registered. Total clients: %d\n", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				fmt.Printf("Client unregistered. Total clients: %d\n", len(h.clients))
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// RegisterClient adds a client to the hub
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// UnregisterClient removes a client from the hub
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// BroadcastMessage sends a message to all connected clients
func (h *Hub) BroadcastMessage(message []byte) {
	h.broadcast <- message
}

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
