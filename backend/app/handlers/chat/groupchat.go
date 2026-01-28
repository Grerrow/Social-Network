package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/app/models"
	"social-network/db"
	"strconv"
	"time"
)

// get the users groups from the group_members table
type GroupConversations struct {
	ID   int    `json:"id"`
	Name string `json:"group_name"`
}

func GetUserGroups(userID int) []GroupConversations {
	var groups []GroupConversations
	rows, err := db.Database.Query(
		`SELECT groups.id, groups.group_name 
         FROM group_members 
         JOIN groups ON group_members.group_id = groups.id
         WHERE group_members.user_id = ?`, userID,
	)

	if err != nil {
		log.Println("Error querying user groups:", err)
		return groups
	}
	defer rows.Close()

	for rows.Next() {
		var g GroupConversations
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			log.Println("Error scanning group row:", err)
			continue
		}
		groups = append(groups, g)
	}
	log.Println("Fetched groups for user", userID, ":", groups)
	return groups
}

// GetGroupMessages retrieves messages for a specific group
func GetGroupMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("ctxUserID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupIDStr := r.URL.Query().Get("group_id")
	if groupIDStr == "" {
		http.Error(w, "group_id is required", http.StatusBadRequest)
		return
	}

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid group_id", http.StatusBadRequest)
		return
	}

	// Check if user is a member of the group
	var isMember bool
	err = db.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?)",
		groupID, userID,
	).Scan(&isMember)

	if err != nil || !isMember {
		http.Error(w, "Not a member of this group", http.StatusForbidden)
		return
	}

	// Fetch messages
	rows, err := db.Database.Query(`
        SELECT gm.id, gm.group_id, gm.sender_id, u.username, gm.content, gm.created_at
        FROM group_messages gm
        JOIN users u ON gm.sender_id = u.id
        WHERE gm.group_id = ?
        ORDER BY gm.created_at ASC
    `, groupID)

	if err != nil {
		log.Printf("[CHAT] Failed to fetch group messages: %v", err)
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	messages := []models.GroupMessage{}
	for rows.Next() {
		var msg models.GroupMessage
		var createdAt time.Time

		err := rows.Scan(&msg.ID, &msg.GroupID, &msg.SenderID, &msg.SenderName, &msg.Content, &createdAt)
		if err != nil {
			log.Printf("[CHAT] Scan error: %v", err)
			continue
		}

		// Convert time.Time to Unix timestamp
		msg.CreatedAt = createdAt.Unix()
		messages = append(messages, msg)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// save group message to database func
func SaveGroupMessage(groupID, senderID, content string, timestamp int64) {
	message, err := db.Database.Prepare("INSERT INTO group_messages (group_id, sender_id, content, timestamp) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println("Error Preparing DB:", err)
		return
	}
	defer message.Close()

	_, err = message.Exec(groupID, senderID, content, timestamp)
	if err != nil {
		log.Println("Error saving group message:", err)
	}
}
