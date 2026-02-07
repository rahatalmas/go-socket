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
	Type         string
	Hub          *Hub
	Conn         *websocket.Conn
	Send         chan []byte
	Customer     *models.Customer
	User         *models.User
	Conversation *models.Conversation
	CancelAI     context.CancelFunc
	SosFlag      bool // -> true when customer talking to human or need to talk to human
	FlagRevealed bool // -> when a human accepts connection
}

// Hub maintains active clients and broadcasts messages
type Hub struct {
	// Registered clients
	clients  map[string]*Client
	allUsers map[string]*Client

	//registered companies to departments to users
	users map[string]map[string]map[*Client]bool

	//customer to user meeting room
	meetingRoom map[*Client]*Client

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
		clients:    make(map[string]*Client),
		allUsers:   make(map[string]*Client),
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
			if client.User != nil || client.Type == "user" {
				h.allUsers[client.User.UserID] = client
				fmt.Printf("user client registered. Total clients: %d\n", len(h.allUsers))
			} else {
				h.clients[client.Customer.Id] = client
				fmt.Printf("customer client registered. Total clients: %d\n", len(h.clients))
			}
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if client.User != nil || client.Type == "user" {
				if _, ok := h.allUsers[client.User.UserID]; ok {
					delete(h.clients, client.User.UserID)
					close(client.Send)
					fmt.Printf("customer client unregistered. Total clients: %d\n", len(h.clients))
				}
			} else {
				if _, ok := h.clients[client.Customer.Id]; ok {
					delete(h.clients, client.Customer.Id)
					close(client.Send)
					fmt.Printf("customer client unregistered. Total clients: %d\n", len(h.clients))
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case h.clients[client].Send <- message:
				default:
					close(h.clients[client].Send)
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

func (h *Hub) GetAllUserCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.allUsers)
}

func (h *Hub) GetAllClients() map[string]*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.clients
}

func (h *Hub) GetAllUsers() map[string]*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.allUsers
}
