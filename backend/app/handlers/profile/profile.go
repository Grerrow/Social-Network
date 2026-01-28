package profile

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"social-network/app/middleware"
	"social-network/app/models"
	"social-network/db"
)

// Profile struct for frontend
type Profile struct {
	models.User
	Posts          []models.Post `json:"posts"`
	Followers      []models.User `json:"followers"`
	Following      []models.User `json:"following"`
	PostsCount     int           `json:"posts_count"`
	FollowersCount int           `json:"followers_count"`
	FollowingCount int           `json:"following_count"`
	Public         bool          `json:"public"`
	Owner          bool          `json:"owner"`
	FollowStatus   string        `json:"follow_status"` // "following", "requested", "none"
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentUserID := middleware.GetUserId(middleware.GetCookieValue(r))
	userProfileID := r.URL.Query().Get("user_id")

	var userID int
	var err error
	if userProfileID == "" {
		userID = currentUserID
	} else {
		userID, err = strconv.Atoi(userProfileID)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
	}
	//check if user exists
	var exists int
	err = db.Database.QueryRow("SELECT 1 FROM users WHERE id = ?", userID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	var profileData Profile

	if userID == currentUserID {
		// Owner: full profile
		profileData = getProfileData(userID, true)
	} else {
		isPublic := getProfilePrivacy(userID)

		if isPublic {
			// Public profile → full profile visible
			profileData = getProfileData(userID, false)
		} else {
			// Private profile → check if current user is a follower
			var isFollower bool
			err := db.Database.QueryRow(
				"SELECT COUNT(*) > 0 FROM followers WHERE followed_id = ? AND follower_id = ?",
				userID, currentUserID,
			).Scan(&isFollower)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}

			if !isFollower {
				// Private profile, not a follower → only basic info + counts
				profileData = getPrivateProfileInfo(userID)
				profileData.FollowStatus = getFollowStatus(userID, currentUserID)
				writeJSON(w, http.StatusOK, profileData)
				return
			}

			// Private profile, follower → full profile visible
			profileData = getProfileData(userID, false)
		}

		// Always set follow status for other users
		profileData.FollowStatus = getFollowStatus(userID, currentUserID)
	}

	writeJSON(w, http.StatusOK, profileData)
}

// ------------------- HELPERS -------------------

func getProfileData(userID int, isOwner bool) Profile {
	followers := GetUserFollowers(userID)
	following := GetUserFollowing(userID)
	posts := getUserPosts(userID)

	return Profile{
		User:      getUserInfo(userID),
		Posts:     posts,
		Followers: followers,
		Following: following,
		Public:    getProfilePrivacy(userID),
		Owner:     isOwner,
	}
}

func getUserInfo(userID int) models.User {
	var user models.User
	err := db.Database.QueryRow(
		`SELECT id, email, first_name, last_name, date_of_birth, avatar, username, about_me, is_private, created_at
		 FROM users WHERE id = ?`, userID,
	).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.DateOfBirth,
		&user.Avatar, &user.Username, &user.AboutMe, &user.IsPrivate, &user.CreatedAt)
	if err != nil {
		fmt.Println("Error scanning user info:", err)
	}
	return user
}

func getUserPosts(userID int) []models.Post {
	rows, err := db.Database.Query("SELECT id, user_id, content, image, created_at FROM posts WHERE user_id = ?", userID)
	if err != nil {
		fmt.Println("Error querying posts:", err)
		return []models.Post{}
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.Image, &post.CreatedAt); err != nil {
			fmt.Println("Error scanning post:", err)
			continue
		}
		posts = append(posts, post)
	}
	return posts
}

func GetUserFollowers(userID int) []models.User {
	rows, err := db.Database.Query("SELECT id, username, avatar FROM users WHERE id IN (SELECT follower_id FROM followers WHERE followed_id = ?)", userID)
	if err != nil {
		fmt.Println("Error querying followers:", err)
		return []models.User{}
	}
	defer rows.Close()

	var followers []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Avatar); err != nil {
			fmt.Println("Error scanning follower:", err)
			continue
		}
		user.UnreadCount = getUnreadCount(userID, user.ID)
		followers = append(followers, user)
	}
	return followers
}

func GetUserFollowing(userID int) []models.User {
	rows, err := db.Database.Query("SELECT id, username, avatar FROM users WHERE id IN (SELECT followed_id FROM followers WHERE follower_id = ?)", userID)
	if err != nil {
		fmt.Println("Error querying following:", err)
		return []models.User{}
	}
	defer rows.Close()

	var following []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Avatar); err != nil {
			fmt.Println("Error scanning following user:", err)
			continue
		}
		user.UnreadCount = getUnreadCount(userID, user.ID)
		following = append(following, user)
	}
	return following
}

func getProfilePrivacy(userID int) bool {
	var isPrivate bool
	err := db.Database.QueryRow("SELECT is_private FROM users WHERE id = ?", userID).Scan(&isPrivate)
	if err != nil {
		fmt.Println("Error fetching privacy:", err)
		return false
	}
	return !isPrivate // true if public
}

func getFollowStatus(profileUserID, currentUserID int) string {
	var exists int
	err := db.Database.QueryRow(
		"SELECT 1 FROM followers WHERE follower_id = ? AND followed_id = ?",
		currentUserID, profileUserID,
	).Scan(&exists)
	if err == nil && exists == 1 {
		return "following"
	}
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("DB error checking followers:", err)
		return "none"
	}

	// Check pending request
	err = db.Database.QueryRow(
		`SELECT 1 FROM follow_user_requests WHERE requester_id = ? AND userToFollow_id = ? AND status='pending'`,
		currentUserID, profileUserID,
	).Scan(&exists)
	if err == nil && exists == 1 {
		return "requested"
	}

	return "none"
}

func getPrivateProfileInfo(userID int) Profile {
	user := getUserInfo(userID)

	var postsCount, followersCount, followingCount int

	// Count posts
	err := db.Database.QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = ?", userID).Scan(&postsCount)
	if err != nil {
		postsCount = 0
	}

	// Count followers
	err = db.Database.QueryRow("SELECT COUNT(*) FROM followers WHERE followed_id = ?", userID).Scan(&followersCount)
	if err != nil {
		followersCount = 0
	}

	// Count following
	err = db.Database.QueryRow("SELECT COUNT(*) FROM followers WHERE follower_id = ?", userID).Scan(&followingCount)
	if err != nil {
		followingCount = 0
	}

	return Profile{
		User: models.User{
			ID:       user.ID,
			Username: user.Username,
			Avatar:   user.Avatar,
			AboutMe:  user.AboutMe,
		},
		Posts:          nil,
		Followers:      nil,
		Following:      nil,
		PostsCount:     postsCount,
		FollowersCount: followersCount,
		FollowingCount: followingCount,
		Public:         false,
		Owner:          false,
	}
}

// the receiver is the current user and the senderID is the user at conversation chatsidebar
func getUnreadCount(receiverID int, senderID int) int {
	var count int
	err := db.Database.QueryRow(
		"SELECT COUNT(*) FROM messages WHERE receiver_id = ? AND sender_id = ? AND is_read = 0",
		receiverID, senderID,
	).Scan(&count)
	if err != nil {
		fmt.Println("Error fetching unread count:", err)
		return 0
	}
	return count
}
