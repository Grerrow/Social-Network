package models

import "time"


type Notification struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	FollowRequestID *int64    `json:"follow_request_id,omitempty"`
	Type            string    `json:"type"`
	SenderID        *int      `json:"sender_id,omitempty"`
	SenderName      *string   `json:"sender_name,omitempty"`
	SenderAvatar    *string   `json:"sender_avatar,omitempty"`
	PostID          *int      `json:"post_id,omitempty"`
	GroupID         *int      `json:"group_id,omitempty"`
	GroupName       *string   `json:"group_name,omitempty"`
	EventID         *int      `json:"event_id,omitempty"`
	EventDate       *int64    `json:"event_date,omitempty"`
	EventTitle      *string   `json:"event_title,omitempty"`
	IsRead          bool      `json:"is_read"`
	CreatedAt       time.Time `json:"created_at"`
}
