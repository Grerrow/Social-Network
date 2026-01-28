package models

type SessionStruct struct {
	ID        string `json:"id"`
	UserID    int    `json:"user_id"`
	ExpiresAt string `json:"expires_at"`
}
