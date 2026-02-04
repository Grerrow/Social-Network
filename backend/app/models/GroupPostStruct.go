package models

type GroupPost struct {
	ID        int            `json:"id"`
	GroupID   int            `json:"group_id"`
	UserID    int            `json:"user_id"`
	Content   string         `json:"content"`
	Image     string         `json:"image,omitempty"`
	CreatedAt int64          `json:"created_at"`
	Author    *User          `json:"author,omitempty"`
	Comments  []GroupComment `json:"comments,omitempty"`
}

type CreatePostRequest struct {
	GroupID string `json:"group_id"`
	Content string `json:"content"`
	Image   string `json:"image,omitempty"`
}
