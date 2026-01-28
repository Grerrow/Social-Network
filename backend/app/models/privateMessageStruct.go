package models

type SendMessageRequest struct {
	ReceiverID int    `json:"receiver_id"`
	Content    string `json:"content"`
}

type Message struct {
	ID         int    `json:"id"`
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"`
	Content    string `json:"content"`
	CreatedAt  int64  `json:"created_at"`
	IsRead     bool   `json:"is_read"`
}
