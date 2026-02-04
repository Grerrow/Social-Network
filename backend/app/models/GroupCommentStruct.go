package models

type GroupComment struct {
	ID        int    `json:"id"`
	PostID    int    `json:"post_id"`
	UserID    int    `json:"user_id"`
	Content   string `json:"content"`
	Image     string `json:"image,omitempty"`
	CreatedAt int64  `json:"created_at"`
	Author    *User  `json:"author,omitempty"`
}

type CreateCommentRequest struct {
	GroupID int    `json:"group_id"`
	PostID  int    `json:"post_id"`
	Content string `json:"content"`
	Image   string `json:"image,omitempty"`
}
