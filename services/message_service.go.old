package services

import (
	"log"
	"net/http"
	"sync"

	"websocket-server/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for simplicity; restrict as needed in production
		return true
	},
}

type Client struct {
	ID     string
	Socket *websocket.Conn
	Send   chan models.Message
}

type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan models.Message
	mu         sync.Mutex
}

var HubInstance = Hub{
	Clients:    make(map[string]*Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Broadcast:  make(chan models.Message),
}

// RunHub starts the hub to handle WebSocket events
func RunHub() {
	for {
		select {
		case client := <-HubInstance.Register:
			HubInstance.mu.Lock()
			HubInstance.Clients[client.ID] = client
			HubInstance.mu.Unlock()
			log.Printf("Client connected: %s", client.ID)

		case client := <-HubInstance.Unregister:
			HubInstance.mu.Lock()
			if _, ok := HubInstance.Clients[client.ID]; ok {
				close(client.Send)
				delete(HubInstance.Clients, client.ID)
				log.Printf("Client disconnected: %s", client.ID)
			}
			HubInstance.mu.Unlock()

		case message := <-HubInstance.Broadcast:
			HubInstance.mu.Lock()
			for _, client := range HubInstance.Clients {
				client.Send <- message
			}
			HubInstance.mu.Unlock()
		}
	}
}

// HandleWebSocket manages WebSocket connections
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	log.Println(r, r.URL, r.URL.Path)

	// Retrieve user ID (example: from query params)
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		log.Println("User ID is required")
		conn.Close()
		return
	}

	client := &Client{
		ID:     userID,
		Socket: conn,
		Send:   make(chan models.Message),
	}

	HubInstance.Register <- client

	// Start goroutines for reading and writing
	go handleRead(client)
	go handleWrite(client)
}

// handleRead reads messages from the WebSocket connection
func handleRead(client *Client) {
	defer func() {
		HubInstance.Unregister <- client
		client.Socket.Close()
	}()

	for {
		var msg models.Message
		err := client.Socket.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// Add the message to the broadcast channel
		HubInstance.Broadcast <- msg
	}
}

// handleWrite writes messages to the WebSocket connection
func handleWrite(client *Client) {
	defer func() {
		client.Socket.Close()
	}()

	for msg := range client.Send {
		err := client.Socket.WriteJSON(msg)
		if err != nil {
			log.Printf("Error writing message: %v", err)
			break
		}
	}
}
