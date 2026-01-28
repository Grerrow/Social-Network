package groups

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"social-network/db"
)

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

func CreateGroupPost(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("ctxUserID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Content == "" && req.Image == "" {
		http.Error(w, "Content or Image is required", http.StatusBadRequest)
		return
	}

	now := time.Now().Unix()
	groupIDint, err := strconv.Atoi(req.GroupID)
	if err != nil || groupIDint == 0 {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}
	// }

	// groupIDInt, _ := strconv.Atoi(groupID)

	// Check if user is a member
	var isMember bool
	db.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?)",
		groupIDint, userID,
	).Scan(&isMember)

	if !isMember {
		http.Error(w, "Must be a member to post", http.StatusForbidden)
		return
	}

	result, err := db.Database.Exec(
		"INSERT INTO group_posts (group_id, user_id, content, image, created_at) VALUES (?, ?, ?, ?, ?)",
		groupIDint, userID, req.Content, req.Image, now,
	)

	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	postID, _ := result.LastInsertId()

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

	post := GroupPost{
		ID:        int(postID),
		GroupID:   groupIDint,
		UserID:    userID,
		Content:   req.Content,
		Image:     req.Image,
		CreatedAt: now,
		Author:    &author,
		Comments:  []GroupComment{},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func ListGroupPosts(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("ctxUserID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID := r.URL.Query().Get("id")

	if groupID == "" {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	// Check if user is a member
	var isMember bool
	db.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?)",
		groupID, userID,
	).Scan(&isMember)

	if !isMember {
		http.Error(w, "Must be a member to view posts", http.StatusForbidden)
		return
	}

	rows, err := db.Database.Query(`
        SELECT gp.id, gp.group_id, gp.user_id, gp.content, gp.image, gp.created_at,
               u.username, u.avatar
        FROM group_posts gp
        JOIN users u ON gp.user_id = u.id
        WHERE gp.group_id = ?
        ORDER BY gp.created_at DESC
    `, groupID)

	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts := []GroupPost{}
	for rows.Next() {
		var p GroupPost
		var author User
		var image, avatar *string

		err := rows.Scan(
			&p.ID, &p.GroupID, &p.UserID, &p.Content, &image, &p.CreatedAt,
			&author.Username, &avatar,
		)
		if err != nil {
			continue
		}

		author.ID = p.UserID
		if avatar != nil {
			author.Avatar = *avatar
		}
		if image != nil {
			p.Image = *image
		}
		p.Author = &author

		// Fetch comments for this post
		p.Comments = fetchCommentsForPost(p.ID)

		posts = append(posts, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func fetchCommentsForPost(postID int) []GroupComment {
	rows, err := db.Database.Query(`
        SELECT gc.id, gc.post_id, gc.user_id, gc.content, gc.image, gc.created_at,
               u.username, u.avatar
        FROM group_post_comments gc
        JOIN users u ON gc.user_id = u.id
        WHERE gc.post_id = ?
        ORDER BY gc.created_at ASC
    `, postID)

	if err != nil {
		return []GroupComment{}
	}
	defer rows.Close()

	comments := []GroupComment{}
	for rows.Next() {
		var c GroupComment
		var author User
		var avatar *string

		err := rows.Scan(
			&c.ID, &c.PostID, &c.UserID, &c.Content, &c.Image, &c.CreatedAt,
			&author.Username, &avatar,
		)
		if err != nil {
			continue
		}

		author.ID = c.UserID
		if avatar != nil {
			author.Avatar = *avatar
		}
		c.Author = &author
		comments = append(comments, c)
	}

	return comments
}
