package models

type GroupMessage struct {
	ID         int    `json:"id"`
	GroupID    int    `json:"group_id"`
	SenderID   int    `json:"sender_id"`
	SenderName string `json:"sender_name,omitempty"`
	Content    string `json:"content"`
	CreatedAt  int64  `json:"created_at"`
}

type SendGroupMessageRequest struct {
	GroupID int    `json:"group_id"`
	Content string `json:"content"`
}
