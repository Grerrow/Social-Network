package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"social-network/app/generalfuncs"
	"social-network/app/handlers/websocket"
	"social-network/app/models"
	"social-network/db"
)

// general function to send message to receiver with appropriate json information
// for the frontend to process
// private chat
func SendMessage(w http.ResponseWriter, r *http.Request) {
	senderID, ok := r.Context().Value("ctxUserID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID := strconv.Itoa(senderID)

	var req models.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Content == "" || req.ReceiverID == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	createdAt := time.Now().Unix()

	result, err := db.Database.Exec(
		"INSERT INTO messages (sender_id, receiver_id, content, created_at, is_read) VALUES (?, ?, ?, ?, 0)",
		senderID, req.ReceiverID, req.Content, createdAt,
	)
	if err != nil {
		log.Printf("[CHAT] DB error: %v", err)
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	msgID, _ := result.LastInsertId()

	msg := models.Message{
		ID:         int(msgID),
		SenderID:   senderID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
		CreatedAt:  createdAt,
		IsRead:     false,
	}

	wsMsg := websocket.WebSocketMessage{
		Type: "private_message",
		Data: msg,
	}

	receiverIDStr := strconv.Itoa(req.ReceiverID)
	log.Printf("[CHAT] Sending message to receiver %s", receiverIDStr)
	websocket.SendToUser(receiverIDStr, wsMsg)
	SendMessageNotification(receiverIDStr, senderID)
	log.Printf("[CHAT] Sending message to sender %s", userID)
	websocket.SendToUser(userID, wsMsg)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

// group chat
func SendGroupMessage(w http.ResponseWriter, r *http.Request) {
	senderID, ok := r.Context().Value("ctxUserID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.SendGroupMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Content == "" || req.GroupID == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	createdAt := time.Now().Unix()

	// Save to database - use createdAt as Unix timestamp
	result, err := db.Database.Exec(
		"INSERT INTO group_messages (group_id, sender_id, content, created_at) VALUES (?, ?, ?, ?)",
		req.GroupID, senderID, req.Content, createdAt,
	)
	if err != nil {
		log.Printf("[CHAT] DB error saving group message: %v", err)
		http.Error(w, "Failed to save group message", http.StatusInternalServerError)
		return
	}

	msgID, _ := result.LastInsertId()

	// Get group members to broadcast to
	memberIDs, err := generalfuncs.GetGroupMemberIDs(req.GroupID, senderID)
	if err != nil {
		log.Printf("[CHAT] Error fetching group members: %v", err)
		http.Error(w, "Failed to fetch group members", http.StatusInternalServerError)
		return
	}

	// Get sender's username for display
	var senderName string
	db.Database.QueryRow("SELECT username FROM users WHERE id = ?", senderID).Scan(&senderName)

	// Build the message
	groupMsg := models.GroupMessage{
		ID:         int(msgID),
		GroupID:    req.GroupID,
		SenderID:   senderID,
		SenderName: senderName,
		Content:    req.Content,
		CreatedAt:  createdAt,
	}

	wsMsg := websocket.WebSocketMessage{
		Type: "group_message",
		Data: groupMsg,
	}

	// Broadcast to all group members (including sender)
	websocket.SendToUsers(memberIDs, wsMsg)
	websocket.SendToUser(strconv.Itoa(senderID), wsMsg)

	log.Printf("[CHAT] Group message sent to group %d by user %d", req.GroupID, senderID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groupMsg)
}

// Send notification for new message
func SendMessageNotification(receiverID string, senderID int) {
	if websocket.IsUserOnline(receiverID) {
		wsMsg := websocket.WebSocketMessage{
			Type: "new_message_notification",
			Data: senderID,
		}
		websocket.SendToUser(receiverID, wsMsg)
	}
}
