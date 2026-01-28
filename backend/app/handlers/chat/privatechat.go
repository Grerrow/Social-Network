package chat

// chat between to clients using websockets

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"social-network/app/handlers/profile"
	"social-network/db"
)

// save message to database func
func SavePrivateMessage(senderID int, receiverID int, content string, timestamp int64) {
	message, err := db.Database.Prepare("INSERT INTO messages (sender_id, receiver_id, content, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return
	}
	defer message.Close()

	_, err = message.Exec(senderID, receiverID, content, time.Unix(timestamp, 0))
	if err != nil {
		return
	}
}

func GetConversations(w http.ResponseWriter, r *http.Request) {
	// get users followers and user that are following him
	// who are online and offline to appear at the conversations list
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("ctxUserID").(int)
	// fill the list with followers or followed users from db
	followers := profile.GetUserFollowers(userID)
	following := profile.GetUserFollowing(userID)
	groups := GetUserGroups(userID)

	// Send the list of followers as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"followers": followers,
		"following": following,
		"groups":    groups,
	})
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	chatId := r.URL.Query().Get("user_id")
	if chatId == "" {
		http.Error(w, "Missing chatId parameter", http.StatusBadRequest)
		return
	}
	userIDint := r.Context().Value("ctxUserID").(int)
	userID := strconv.Itoa(userIDint)
	// Query messages between current user and chatId/receiverID user
	rows, err := db.Database.Query(`
        SELECT sender_id, receiver_id, content, created_at 
        FROM messages 
        WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
        ORDER BY created_at ASC
    `, userID, chatId, chatId, userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []map[string]interface{}
	for rows.Next() {
		var senderID, receiverID int
		var content string
		var createdAt time.Time

		if err := rows.Scan(&senderID, &receiverID, &content, &createdAt); err != nil {
			continue
		}

		messages = append(messages, map[string]interface{}{
			"sender_id":   senderID,
			"receiver_id": receiverID,
			"content":     content,
			"created_at":  createdAt.Unix(),
		})
	}

	if messages == nil {
		messages = []map[string]interface{}{}
	}
	chatID, _ := strconv.Atoi(chatId)
	// mark messages as read
	MarkMessagesAsRead(chatID, userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// mark messages as read
func MarkMessagesAsRead(senderID int, receiverID string) error {
	_, err := db.Database.Exec(
		"UPDATE messages SET is_read = 1 WHERE sender_id = ? AND receiver_id = ? AND is_read = 0",
		senderID, receiverID,
	)
	return err
}
