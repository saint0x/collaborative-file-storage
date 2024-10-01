package websocket

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/saint0x/file-storage-app/backend/internal/db"
	"github.com/saint0x/file-storage-app/backend/internal/models"
)

// UpdateType represents different types of updates that can be sent
type UpdateType string

const (
	FileUploaded      UpdateType = "file_uploaded"
	FileDeleted       UpdateType = "file_deleted"
	CollectionCreated UpdateType = "collection_created"
	CollectionUpdated UpdateType = "collection_updated"
	CollectionDeleted UpdateType = "collection_deleted"
	UserJoined        UpdateType = "user_joined"
	UserLeft          UpdateType = "user_left"
)

// Update represents a structured update message
type Update struct {
	Type UpdateType  `json:"type"`
	Data interface{} `json:"data"`
}

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan Update

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Mutex for thread-safe operations on the clients map
	mu sync.RWMutex

	// SQLite client for database operations
	db *db.SQLiteClient
}

// NewHub creates a new Hub instance
func NewHub(db *db.SQLiteClient) *Hub {
	return &Hub{
		broadcast:  make(chan Update),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		db:         db,
	}
}

// Run starts the hub and handles client connections and messages
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			h.BroadcastUpdate(UserJoined, client.ID)
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.BroadcastUpdate(UserLeft, client.ID)
			}
			h.mu.Unlock()
		case update := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- update:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastUpdate sends an update to all connected clients
func (h *Hub) BroadcastUpdate(updateType UpdateType, data interface{}) {
	update := Update{
		Type: updateType,
		Data: data,
	}
	h.broadcast <- update
}

// SendUpdateToClient sends an update to a specific client
func (h *Hub) SendUpdateToClient(client *Client, updateType UpdateType, data interface{}) {
	update := Update{
		Type: updateType,
		Data: data,
	}
	select {
	case client.send <- update:
	default:
		h.mu.Lock()
		close(client.send)
		delete(h.clients, client)
		h.mu.Unlock()
	}
}

// recordPing records a ping in the SQLite database
func (h *Hub) recordPing(clientID string) error {
	ping := models.Ping{
		ID:        uuid.New(),
		ClientID:  clientID,
		Timestamp: time.Now(),
	}

	_, err := h.db.DB.Exec("INSERT INTO pings (id, client_id, timestamp) VALUES (?, ?, ?)",
		ping.ID, ping.ClientID, ping.Timestamp)
	return err
}

// Add this method to the Hub struct
func (h *Hub) Stop() {
	close(h.broadcast)
	// Add any additional cleanup logic here
}
