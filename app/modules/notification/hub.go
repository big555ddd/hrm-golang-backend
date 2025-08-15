package notification

import (
	"app/internal/logger"
	"encoding/json"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients with their user IDs
	Clients map[*Client]string

	// Inbound messages from the clients
	Broadcast chan []byte

	// Register requests from the clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client

	// Mutex to protect concurrent access
	mutex sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]string),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mutex.Lock()
			h.Clients[client] = client.UserID
			h.mutex.Unlock()

			logger.Infof("User %s connected, total clients: %d", client.UserID, len(h.Clients))

		case client := <-h.Unregister:
			h.mutex.Lock()
			if userID, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				logger.Infof("User %s disconnected, remaining clients: %d", userID, len(h.Clients))
			}
			h.mutex.Unlock()

		case message := <-h.Broadcast:
			h.mutex.RLock()
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

// BroadcastToAll sends message to all connected clients
func (h *Hub) BroadcastToAll(message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		logger.Infof("Error marshaling message: %v", err)
		return
	}

	h.Broadcast <- data
}

// SendToUser sends message to specific user's connections
func (h *Hub) SendToUser(userID string, message interface{}) {
	data, err := json.Marshal(message)
	if err != nil {
		logger.Infof("Error marshaling message: %v", err)
		return
	}

	h.mutex.RLock()
	defer h.mutex.RUnlock()

	sentCount := 0
	totalClients := 0

	// Find all clients for this user and send message
	for client, clientUserID := range h.Clients {
		if clientUserID == userID {
			totalClients++
			select {
			case client.Send <- data:
				sentCount++
			default:
				logger.Infof("Client channel full for user: %s", userID)
			}
		}
	}

	if totalClients == 0 {
		logger.Infof("No active connections for user: %s", userID)
	} else {
		logger.Infof("Message sent to %d/%d connections for user: %s", sentCount, totalClients, userID)
	}
} // GetActiveUsers returns list of active user IDs
func (h *Hub) GetActiveUsers() []string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	userSet := make(map[string]bool)
	for _, userID := range h.Clients {
		if userID != "" {
			userSet[userID] = true
		}
	}

	users := make([]string, 0, len(userSet))
	for userID := range userSet {
		users = append(users, userID)
	}

	return users
}

// GetUserConnectionCount returns number of connections for a user
func (h *Hub) GetUserConnectionCount(userID string) int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	count := 0
	for _, clientUserID := range h.Clients {
		if clientUserID == userID {
			count++
		}
	}
	return count
}

// IsUserOnline checks if user has active connections
func (h *Hub) IsUserOnline(userID string) bool {
	return h.GetUserConnectionCount(userID) > 0
}

// GetConnectionInfo returns debug information about connections
func (h *Hub) GetConnectionInfo() map[string]interface{} {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	userCounts := make(map[string]int)
	for _, userID := range h.Clients {
		if userID != "" {
			userCounts[userID]++
		}
	}

	return map[string]interface{}{
		"total_clients": len(h.Clients),
		"user_counts":   userCounts,
		"active_users":  h.GetActiveUsers(),
	}
}
