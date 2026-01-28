package models

import (
	"time"
)

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	// just the date not time
	DateOfBirth string    `json:"date_of_birth"`
	Avatar      *string   `json:"avatar,omitempty"`
	Username    *string   `json:"username,omitempty"`
	AboutMe     *string   `json:"about_me,omitempty"`
	IsPrivate   bool      `json:"is_private"` // Profile privacy: false = public, true = private
	CreatedAt   time.Time `json:"created_at"`
	UnreadCount int       `json:"unread_count,omitempty"`
}

// Avatar, Username, AboutMe are pointers (*string) to allow NULL values in the database
// omitempty: JSON tag: excludes these fields from JSON response when nil
