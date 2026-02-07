package handler

import (
	"butter-socket/internal/hub"
	"butter-socket/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512 * 1024 // 512KB
)

// WsHandler handles WebSocket connections
func WsHandler(h *hub.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error while upgrading connection:", err)
		return
	}

	// Extract query parameters
	queryParams := r.URL.Query()
	customerId := queryParams.Get("customer_id")
	companyId := queryParams.Get("company_id")
	source := queryParams.Get("source")

	// Validate required parameters
	if customerId == "" {
		log.Println("Missing required parameter: customer_id")
		errMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Missing customer_id parameter")
		conn.WriteMessage(websocket.CloseMessage, errMsg)
		conn.Close()
		return
	}

	if companyId == "" {
		log.Println("Missing required parameter: company_id")
		errMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Missing company_id parameter")
		conn.WriteMessage(websocket.CloseMessage, errMsg)
		conn.Close()
		return
	}

	// If source is not provided in query params, fall back to Origin header
	if source == "" {
		source = r.Header.Get("Origin")
		if source == "" {
			source = "unknown"
		}
	}

	// Create a new client with parameters from query string
	client := &hub.Client{
		Hub:  h,
		Conn: conn,
		Send: make(chan []byte, 256),
		Customer: &models.Customer{
			Id:        customerId,
			CompanyId: companyId,
			Source:    source,
		},
		Conversation: &models.Conversation{
			Id:        uuid.New().String(),
			CompanyId: companyId,
			Status:    "open",
			Customer: &models.Customer{
				Id:        customerId,
				CompanyId: companyId,
				Source:    source,
			},
		},
	}

	log.Printf("New WebSocket connection - Customer ID: %s, Company ID: %s, Source: %s",
		customerId, companyId, source)

	// Register the client
	client.Hub.RegisterClient(client)

	// Send welcome message
	sendWelcomeMessage(client)

	// Start goroutines for reading and writing
	go writePump(client)
	go readPump(client)
}

// readPump reads messages from the WebSocket connection
func readPump(client *hub.Client) {
	defer func() {
		client.Hub.UnregisterClient(client)
		client.Conn.Close()
	}()

	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetReadLimit(maxMessageSize)
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Process the incoming message
		handleIncomingMessage(client, message)
	}
}

// writePump writes messages to the WebSocket connection
func writePump(client *hub.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(client.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-client.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleIncomingMessage processes incoming messages from clients
func handleIncomingMessage(client *hub.Client, message []byte) {
	var wsMsg models.WSMessage
	if err := json.Unmarshal(message, &wsMsg); err != nil {
		sendError(client, "Invalid WS message format")
		return
	}

	switch wsMsg.Type {
	case "transfer_chat":
		handleChatTransferToUser(client, wsMsg.Payload)
	case "accept_chat":
		handleHumanAcceptTheChat(client, wsMsg.Payload)
	case "message":
		fmt.Println("client id: ", client.Customer.Id)
		fmt.Println("flags: ", client.SosFlag, client.FlagRevealed)
		if client.FlagRevealed == true {
			handleConversationWithHuman(client, wsMsg.Payload)
		} else {
			handleChatStreamMessage(client, wsMsg.Payload)
		}
	case "ping":
		sendPong(client)
	default:
		sendError(client, "Unknown message type")
	}
}

// sendWelcomeMessage sends a welcome message to newly connected clients
func sendWelcomeMessage(client *hub.Client) {
	// systemMsg := models.Message{
	// 	MetaData: models.MetaData{
	// 		CreatedAt: time.Now().Format(time.RFC3339),
	// 	},
	// 	Id:             uuid.New().String(),
	// 	ConversationId: client.Conversation.Id,
	// 	SenderId:       "system",
	// 	SenderType:     "system",
	// 	Content: fmt.Sprintf(
	// 		"Welcome! Customer ID: %s, Company ID: %s",
	// 		client.Customer.Id,
	// 		client.Customer.CompanyId,
	// 	),
	// 	ContentType: "text",
	// }
	var msgOut = models.MsgInOut{
		SenderType:  "AI-AGENT",
		Content:     "Welcome to Butter Chat",
		ContentType: "txt",
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	sendMessage(client, "welcome", msgOut)
}

// sendMessage sends a message to a specific client
func sendMessage(client *hub.Client, msgType string, payload interface{}) {
	wsMsg := models.WSMessage{
		Type:    msgType,
		Payload: payload,
	}

	msgBytes, err := json.Marshal(wsMsg)
	if err != nil {
		log.Println("Error marshaling message:", err)
		return
	}

	select {
	case client.Send <- msgBytes:
	default:
		log.Println("Client send channel is full")
	}
}

// sendError sends an error message to the client
func sendError(client *hub.Client, errorMsg string) {
	errorPayload := map[string]string{
		"error": errorMsg,
	}
	sendMessage(client, "error", errorPayload)
}

// sendPong responds to ping messages
func sendPong(client *hub.Client) {
	pongPayload := map[string]string{
		"status": "pong",
	}
	sendMessage(client, "pong", pongPayload)
}
