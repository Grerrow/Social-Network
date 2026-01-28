package models

import "time"

type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	GroupID     int       `json:"group_id"`
	CreatorID   int       `json:"creator_id"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
}

type EventResponse struct {
	EventID     int       `json:"event_id"`
	UserID      int       `json:"user_id"`
	Response    string    `json:"response"` // "going" or "not_going"
	RespondedAt time.Time `json:"responded_at"`
}
