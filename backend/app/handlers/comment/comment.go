package comment

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"social-network/app/generalfuncs"
	"social-network/app/middleware"
	"social-network/app/models"
	"social-network/db"
)

func CommentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetComments(w, r)
	case http.MethodPost:
		CreateComment(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var comment models.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Input validation ---------------------------------------------------------------------------
	if comment.Image == nil && strings.TrimSpace(comment.Content) == "" {
		http.Error(w, "Image or text is required", http.StatusBadRequest)
		return
	}
	if len(comment.Content) > 400 {
		http.Error(w, "Content is too long (max 400 characters)", http.StatusBadRequest)
		return
	}
	if comment.Image != nil && *comment.Image != "" && !strings.HasPrefix(*comment.Image, "/images") {
		http.Error(w, "Invalid image path", http.StatusBadRequest)
		return
	}

	// Checking if user is authorized (cookie.Value is session ID) --------------------------------
	cookieValue := middleware.GetCookieValue(r)
	userID := middleware.GetUserId(cookieValue)
	if userID <= 0 {
		http.Error(w, "Not an authorized user", http.StatusBadRequest)
		return
	}

	// set the authenticated user as the comment creator (for security reasons) -------------------
	comment.UserID = userID

	// Insert comment into database
	query := `INSERT INTO comments (content, image, post_id, user_id, created_at) 
VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)`

	result, err := db.Database.Exec(query,
		comment.Content,
		comment.Image,
		comment.PostID,
		comment.UserID,
	)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}
	// insert the notification of the comment for the post owner
	//
	var postOwnerID int
	err = db.Database.QueryRow(
		"SELECT user_id FROM posts WHERE id = ?",
		comment.PostID,
	).Scan(&postOwnerID)
	if err != nil {
		http.Error(w, "Failed to get post owner", http.StatusInternalServerError)
		return
	}
	// avoid notifying oneself
	//
	// Get the ID of the newly created comment
	commentID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to get comment ID", http.StatusInternalServerError)
		return
	}
	comment.ID = int(commentID)

	// Fetch the comment including username, avatar, and created_at to match expected structure
	row := db.Database.QueryRow(`
		SELECT c.id, c.content, c.image, c.post_id, c.user_id, c.created_at, 
		       u.username, u.first_name, u.last_name, u.avatar
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = ?`, comment.ID)

	var username, first_name, last_name, avatar string
	if err := row.Scan(&comment.ID, &comment.Content, &comment.Image, &comment.PostID,
		&comment.UserID, &comment.CreatedAt, &username, &first_name, &last_name, &avatar); err != nil {
		http.Error(w, "Failed to fetch comment", http.StatusInternalServerError)
		return
	}
	comment.Username = username
	comment.FirstName = first_name
	comment.LastName = last_name
	comment.Avatar = &avatar

	if postOwnerID != comment.UserID {
		// insert notification into database
		err = generalfuncs.CreateNotification(models.Notification{
			UserID:     postOwnerID,     // who receives it
			Type:       "new_comment",   // notification type
			SenderID:   &comment.UserID, // who triggered it
			PostID:     &comment.PostID, // related post
			SenderName: &first_name,
		})
		if err != nil {
			http.Error(w, "Failed to create notification", http.StatusInternalServerError)
			return
		}
	}

	// Return created comment as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

// GET /comments?post_id=123
func GetComments(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("post_id")
	if postIDStr == "" {
		http.Error(w, "post_id is required", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "invalid post_id", http.StatusBadRequest)
		return
	}

	query := `
        SELECT c.id, c.content, c.image, c.post_id, c.user_id, c.created_at, u.username, u.first_name, u.last_name, u.avatar
        FROM comments c
        JOIN users u ON c.user_id = u.id
        WHERE c.post_id = ?
        ORDER BY c.created_at ASC
    `
	rows, err := db.Database.Query(query, postID)
	if err != nil {
		http.Error(w, "failed to fetch comments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	comments := []models.Comment{}
	for rows.Next() {
		var c models.Comment
		var username, firstName, lastName, avatar string
		if err := rows.Scan(&c.ID, &c.Content, &c.Image, &c.PostID, &c.UserID, &c.CreatedAt, &username, &firstName, &lastName, &avatar); err != nil {
			http.Error(w, "failed to parse comment", http.StatusInternalServerError)
			return
		}
		c.Username = username
		c.FirstName = firstName
		c.LastName = lastName
		c.Avatar = &avatar
		comments = append(comments, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
