package models

import "time"

type Post struct {
	ID               int       `json:"id"`
	Content          string    `json:"content,omitempty"`
	Image            *string   `json:"image,omitempty"`
	Privacy          string    `json:"privacy"`
	UserID           int       `json:"user_id"`
	GroupID          *int      `json:"group_id,omitempty"` // nil for regular posts or set for group posts
	CreatedAt        time.Time `json:"created_at"`
	Username         string    `json:"username"`
	Avatar           string    `json:"avatar"`
	AllowedFollowers []int     `json:"allowed_followers,omitempty"` // for almost_private posts
}

// Image and GroupID are pointers to allow NULL values in the database
// The 'omitempty' JSON tag excludes these fields from JSON response when nil

// privacy constants are declared here for safety and refactoring safety:
const (
	PrivacyPublic        = "public"
	PrivacyAlmostPrivate = "almost_private"
	PrivacyPrivate       = "private"
)
