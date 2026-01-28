package groups

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"social-network/db"
)

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

func CreateGroupComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("ctxUserID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if user is a member
	var isMember bool
	db.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?)",
		req.GroupID, userID,
	).Scan(&isMember)

	if !isMember {
		http.Error(w, "Must be a member to comment", http.StatusForbidden)
		return
	}

	// Check if post exists and belongs to this group
	var postExists bool
	db.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM group_posts WHERE id = ? AND group_id = ?)",
		req.PostID, req.GroupID,
	).Scan(&postExists)

	if !postExists {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if req.Content == "" && req.Image == "" {
		http.Error(w, "Content or Image is required", http.StatusBadRequest)
		return
	}

	now := time.Now().Unix()

	result, err := db.Database.Exec(
		"INSERT INTO group_post_comments (post_id, user_id, content, image, created_at) VALUES (?, ?, ?, ?, ?)",
		req.PostID, userID, req.Content, req.Image, now,
	)

	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	commentID, _ := result.LastInsertId()

	// Get author info
	var author User
	var avatar *string
	db.Database.QueryRow(
		"SELECT id, username, avatar FROM users WHERE id = ?",
		userID,
	).Scan(&author.ID, &author.Username, &avatar)

	if avatar != nil {
		author.Avatar = *avatar
	}

	comment := GroupComment{
		ID:        int(commentID),
		PostID:    req.PostID,
		UserID:    userID,
		Content:   req.Content,
		Image:     req.Image,
		CreatedAt: now,
		Author:    &author,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func ListGroupComments(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("ctxUserID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID := r.URL.Query().Get("group_id")
	postID := r.URL.Query().Get("post_id")

	if groupID == "" || postID == "" {
		http.Error(w, "Group ID and Post ID are required", http.StatusBadRequest)
		return
	}

	// Check if user is a member
	var isMember bool
	db.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?)",
		groupID, userID,
	).Scan(&isMember)

	if !isMember {
		http.Error(w, "Must be a member to view comments", http.StatusForbidden)
		return
	}

	comments := fetchCommentsForPost(func() int {
		id, _ := strconv.Atoi(postID)
		return id
	}())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
