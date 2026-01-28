package websocket

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type WebSocketClient struct {
	UserID string
	Conn   *websocket.Conn
	send   chan WebSocketMessage
	mu     sync.Mutex
}

type WebSocketHub struct {
	clients    map[*WebSocketClient]bool
	userIndex  map[string]map[*WebSocketClient]bool
	register   chan *WebSocketClient
	unregister chan *WebSocketClient
	mu         sync.RWMutex
}

var Hub = &WebSocketHub{
	clients:    make(map[*WebSocketClient]bool),
	userIndex:  make(map[string]map[*WebSocketClient]bool),
	register:   make(chan *WebSocketClient),
	unregister: make(chan *WebSocketClient),
}

func StartHub() {
	go Hub.run()
}

func (h *WebSocketHub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			if h.userIndex[client.UserID] == nil {
				h.userIndex[client.UserID] = make(map[*WebSocketClient]bool)
			}
			h.userIndex[client.UserID][client] = true
			count := len(h.userIndex[client.UserID])
			h.mu.Unlock()

			log.Printf("[WS] User %s connected (connections: %d)", client.UserID, count)

			// Only broadcast presence (online/offline) - NOT private messages
			h.broadcastPresence(client.UserID, true)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			if sockets, ok := h.userIndex[client.UserID]; ok {
				delete(sockets, client)
				if len(sockets) == 0 {
					delete(h.userIndex, client.UserID)
					h.mu.Unlock()
					// Broadcast that user went offline
					h.broadcastPresence(client.UserID, false)
					log.Printf("[WS] User %s disconnected", client.UserID)
				} else {
					h.mu.Unlock()
				}
			} else {
				h.mu.Unlock()
			}
		}
	}
}

// broadcastPresence - sends online/offline status to ALL users
// This is fine because everyone should see who is online
// optional improvement for scalability
// only broadcast to friends/followers
func (h *WebSocketHub) broadcastPresence(userID string, online bool) {
	msg := WebSocketMessage{
		Type: "presence",
		Data: map[string]interface{}{
			"user_id": userID,
			"online":  online,
		},
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		select {
		case client.send <- msg:
		default:
		}
	}
}

// broadcast to only one user
func (h *WebSocketHub) SendToUser(userID string, msg WebSocketMessage) {
	h.mu.RLock()
	clients := h.userIndex[userID]
	h.mu.RUnlock()

	if clients == nil {
		log.Printf("[WS] User %s not connected", userID)
		return
	}

	for client := range clients {
		select {
		case client.send <- msg:
			log.Printf("[WS] Sent %s to user %s", msg.Type, userID)
		default:
			log.Printf("[WS] Channel full for user %s", userID)
		}
	}
}

// send to multiple users from group members list
func SendToUsers(userIDs []int, msg WebSocketMessage) {
	for _, userID := range userIDs {
		Hub.SendToUser(strconv.Itoa(userID), msg)
	}
}

// GetOnlineUsers - returns list of online user IDs
func (h *WebSocketHub) GetOnlineUsers() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]string, 0, len(h.userIndex))
	for userID := range h.userIndex {
		users = append(users, userID)
	}
	return users
}

// Exported function for other packages to send private messages
func SendToUser(userID string, msg WebSocketMessage) {
	Hub.SendToUser(userID, msg)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	userIDint, ok := r.Context().Value("ctxUserID").(int)
	if !ok || userIDint == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := strconv.Itoa(userIDint)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] Upgrade error: %v", err)
		return
	}

	client := &WebSocketClient{
		UserID: userID,
		Conn:   conn,
		send:   make(chan WebSocketMessage, 256),
	}

	Hub.register <- client

	// Send initial online users list to THIS user only
	go func() {
		time.Sleep(100 * time.Millisecond)
		onlineUsers := Hub.GetOnlineUsers()
		client.send <- WebSocketMessage{
			Type: "presence:init",
			Data: onlineUsers,
		}
	}()

	go client.writePump()
	go client.readPump()
}

func (c *WebSocketClient) readPump() {
	defer func() {
		Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(8192)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg WebSocketMessage
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WS] Read error for user %s: %v", c.UserID, err)
			}
			break
		}

		// Handle ping from client
		if msg.Type == "ping" {
			c.send <- WebSocketMessage{Type: "pong"}
		}
	}
}

func (c *WebSocketClient) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.mu.Lock()
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			err := c.Conn.WriteJSON(msg)
			c.mu.Unlock()

			if err != nil {
				log.Printf("[WS] Write error for user %s: %v", c.UserID, err)
				return
			}

		case <-ticker.C:
			c.mu.Lock()
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			err := c.Conn.WriteMessage(websocket.PingMessage, nil)
			c.mu.Unlock()

			if err != nil {
				return
			}
		}
	}
}

// check if user that receives messages is online
func IsUserOnline(userID string) bool {
	Hub.mu.RLock()
	defer Hub.mu.RUnlock()

	clients, ok := Hub.userIndex[userID]
	return ok && len(clients) > 0
}
