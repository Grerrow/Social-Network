package models

import "time"

type Group struct {
	ID          int       `json:"id"`
	GroupName   string    `json:"group_name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatorID   int       `json:"creator_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type GroupMember struct {
	ID      int `json:"id"`
	GroupID int `json:"group_id"`
	UserID  int `json:"user_id"`
}
