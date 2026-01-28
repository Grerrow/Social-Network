package models

import "time"

type Comment struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Image     *string   `json:"image,omitempty"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Avatar    *string   `json:"avatar,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Image is a pointer to allow NULL values in the database
// The 'omitempty' JSON tag excludes this field from JSON response when nil
