package post

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"social-network/app/middleware"
	"social-network/app/models"
	"social-network/db"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var post models.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// extra validation ----------------------------------------------------------------
	if post.Image == nil && strings.TrimSpace(post.Content) == "" {
		http.Error(w, "Image or text is required", http.StatusBadRequest)
		return
	}
	if len(post.Content) > 400 {
		http.Error(w, "Content is too long (max 400 characters)", http.StatusBadRequest)
		return
	}
	if post.Image != nil && *post.Image != "" && !strings.HasPrefix(*post.Image, "/images") {
		http.Error(w, "Invalid image path", http.StatusBadRequest)
		return
	}

	cookieValue := middleware.GetCookieValue(r)
	userID := middleware.GetUserId(cookieValue)
	if userID <= 0 {
		http.Error(w, "Not an authorized user", http.StatusBadRequest)
		return
	}

	// set the authenticated user as the post creator (for security reasons) ----------------------
	post.UserID = userID

	// Privacy value validation -------------------------------------------------------------------
	if post.Privacy != models.PrivacyPublic &&
		post.Privacy != models.PrivacyAlmostPrivate &&
		post.Privacy != models.PrivacyPrivate {
		http.Error(w, "Invalid privacy value", http.StatusBadRequest)
		return
	}

	// Insert post into database ------------------------------------------------------------------
	query := `INSERT INTO posts (content, image, privacy, user_id, group_id, created_at) 
	VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`

	result, err := db.Database.Exec(query,
		post.Content,
		post.Image,
		post.Privacy,
		post.UserID,
		post.GroupID,
	)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	// get the ID of the newly created post -------------------------------------------------------
	// i like this lastinsertid function
	postID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to get post ID", http.StatusInternalServerError)
		return
	}
	post.ID = int(postID)

	// Handle almost_private posts - insert allowed followers into post_visibility table
	if post.Privacy == models.PrivacyAlmostPrivate && len(post.AllowedFollowers) > 0 {
		for _, followerID := range post.AllowedFollowers {
			_, err := db.Database.Exec(
				`INSERT INTO post_visibility (post_id, user_id) VALUES (?, ?)`,
				post.ID, followerID,
			)
			if err != nil {
				log.Printf("Failed to insert post visibility for user %d: %v", followerID, err)
				// Continue with other followers even if one fails
			}
		}
	}

	// Return created post as JSON in the form the front end expects -------------------------------
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"post": post,
	})
}

// get posts of people followed by the user, posts from public profile and user's own posts
func GetFeedPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookieValue := middleware.GetCookieValue(r)
	userID := middleware.GetUserId(cookieValue)
	if userID <= 0 {
		http.Error(w, "Not an authorized user", http.StatusBadRequest)
		return
	}

	// Query for posts:
	// - Public posts from followed users or public profiles
	// - Almost private posts where current user is in the allowed list (post_visibility table)
	// - Private posts only from the user themselves
	// - User's own posts regardless of privacy
	query := `
	SELECT DISTINCT p.id, p.content, p.image, p.privacy, p.user_id, p.created_at, u.username, u.avatar
	FROM posts p
	JOIN users u ON p.user_id = u.id
	LEFT JOIN followers f ON p.user_id = f.followed_id AND f.follower_id = ?
	LEFT JOIN post_visibility pv ON p.id = pv.post_id AND pv.user_id = ?
	WHERE 
	    p.group_id IS NULL AND (
	        p.privacy = 'public'
	        OR (p.privacy = 'almost_private' AND pv.user_id IS NOT NULL)
	        OR (p.privacy = 'private' AND f.follower_id IS NOT NULL)
	        OR (p.user_id = ?)
	    )
	ORDER BY p.created_at DESC
	LIMIT 20;`

	rows, err := db.Database.Query(query, userID, userID, userID)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		log.Println("Error querying posts:", err)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var username, avatar string
		if err := rows.Scan(&post.ID, &post.Content, &post.Image, &post.Privacy, &post.UserID, &post.CreatedAt, &username, &avatar); err != nil {
			http.Error(w, "Failed to parse posts", http.StatusInternalServerError)
			log.Println("Error scanning post row:", err)
			return
		}
		post.Username = username
		post.Avatar = avatar
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"posts": posts,
	})
}

// get post data handler - for post view page and comments
func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	postIDStr := r.URL.Query().Get("post_id")
	if postIDStr == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}
	//check if postID exists and is valid integer
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	var post models.Post

	query := `
	SELECT p.id, p.content, p.image, p.privacy, p.user_id, p.created_at, u.username, u.avatar
	FROM posts p
	JOIN users u ON p.user_id = u.id
	WHERE p.id = ?;`

	err = db.Database.QueryRow(query, postID).Scan(&post.ID, &post.Content, &post.Image, &post.Privacy, &post.UserID, &post.CreatedAt, &post.Username, &post.Avatar)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		log.Println("Error querying post:", err)
		return
	}

	//query comments for the post
	commentsQuery := `
	SELECT c.id, c.content, c.image, c.user_id, c.created_at, u.username, u.avatar
	FROM comments c
	JOIN users u ON c.user_id = u.id
	WHERE c.post_id = ?
	ORDER BY c.created_at ASC;`

	rows, err := db.Database.Query(commentsQuery, postID)
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		log.Println("Error querying comments:", err)
		return
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.Image, &comment.UserID, &comment.CreatedAt, &comment.Username, &comment.Avatar); err != nil {
			http.Error(w, "Failed to parse comments", http.StatusInternalServerError)
			log.Println("Error scanning comment row:", err)
			return
		}
		comments = append(comments, comment)
	}

	// Return post with comments as JSON

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"post":     post,
		"comments": comments,
	})
}

// for private post options - get followers of the user
func GetFollowersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("ctxUserID").(int)
	if userID == 0 {
		http.Error(w, "Not an authorized user", http.StatusBadRequest)
		return
	}

	var followers []Followers

	query := `
	SELECT u.id, u.username, u.first_name, u.last_name
	FROM users u
	JOIN followers f ON u.id = f.follower_id
	WHERE f.followed_id = ?;`

	rows, err := db.Database.Query(query, userID)
	if err != nil {
		http.Error(w, "Failed to fetch followers", http.StatusInternalServerError)
		log.Println("Error querying followers:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var follower Followers
		if err := rows.Scan(&follower.ID, &follower.Username, &follower.FirstName, &follower.LastName); err != nil {
			http.Error(w, "Failed to parse followers", http.StatusInternalServerError)
			log.Println("Error scanning follower row:", err)
			return
		}
		followers = append(followers, follower)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"followers": followers,
	})
}

type Followers struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
