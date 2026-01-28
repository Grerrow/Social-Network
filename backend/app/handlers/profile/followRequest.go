package profile

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"social-network/app/generalfuncs"
	"social-network/app/middleware"
	"social-network/app/models"
	"social-network/db"
)

//CREATE TABLE followers (
//    follower_id INTEGER NOT NULL,
//    followed_id INTEGER NOT NULL,
//    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//
//    -- composite primary key (2 primary keys: follower_id, followed_id) ensures each user can only follow another user once
//    PRIMARY KEY (follower_id, followed_id),
//
//    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
//    FOREIGN KEY (followed_id) REFERENCES users(id) ON DELETE CASCADE
//);

// CREATE TABLE follow_user_requests (
//     id INTEGER PRIMARY KEY AUTOINCREMENT,
//     userToFollow_id INTEGER NOT NULL,
//     requester_id INTEGER NOT NULL,
//     status TEXT NOT NULL DEFAULT 'pending', -- 'pending', 'approved', 'rejected'
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

//     -- prevent duplicate follow requests
//     UNIQUE (userToFollow_id, requester_id),

//     FOREIGN KEY (userToFollow_id) REFERENCES users(id) ON DELETE CASCADE,
//     FOREIGN KEY (requester_id) REFERENCES users(id)
// );

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

type FollowRequest struct {
	ID           int       `json:"id"`
	RequesterID  int       `json:"requester_id"`
	SenderName   string    `json:"sender_name,omitempty"`
	SenderAvatar *string   `json:"sender_avatar,omitempty"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

// /follow request to a user and check if public or private profile
// FollowRequestHandler handles follow requests and auto-follow for public profiles
func FollowRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	currentUserID := middleware.GetUserId(middleware.GetCookieValue(r))
	if currentUserID == 0 {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	// get sednders name and sender avatar
	var senderName, senderAvatar string
	err := db.Database.QueryRow(
		"SELECT first_name, avatar FROM users WHERE id = ?",
		currentUserID,
	).Scan(&senderName, &senderAvatar)
	if err != nil {
		log.Println("DB error reading sender name or avatar:", err)
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}

	userToFollowIDStr := r.URL.Query().Get("user_id")
	userToFollowID, err := strconv.Atoi(userToFollowIDStr)
	if err != nil || userToFollowID == 0 {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if userToFollowID == currentUserID {
		writeError(w, http.StatusBadRequest, "cannot follow yourself")
		return
	}

	// Check profile privacy
	var isPrivate bool
	err = db.Database.QueryRow(
		"SELECT is_private FROM users WHERE id = ?",
		userToFollowID,
	).Scan(&isPrivate)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "target user not found")
			return
		}
		log.Println("DB error reading profile public:", err)
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}

	// Check if already following
	var exists int
	err = db.Database.QueryRow(
		"SELECT 1 FROM followers WHERE follower_id = ? AND followed_id = ?",
		currentUserID, userToFollowID,
	).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Println("DB error checking existing follower:", err)
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}
	if exists == 1 {
		writeError(w, http.StatusConflict, "already following")
		return
	}

	// Public profile → auto-follow
	if !isPrivate {
		_, err := db.Database.Exec(
			"INSERT INTO followers (follower_id, followed_id) VALUES (?, ?)",
			currentUserID, userToFollowID,
		)
		if err != nil {
			log.Println("DB error inserting follower:", err)
			writeError(w, http.StatusInternalServerError, "database error")
		}
		// notification for a new follower
		err = generalfuncs.CreateNotification(models.Notification{
			UserID:       userToFollowID,
			Type:         "new_follower",
			SenderID:     &currentUserID,
			SenderName:   &senderName,
			SenderAvatar: &senderAvatar,
		})
		if err != nil {
			log.Println("DB error inserting new follower notification:", err)
		}

		writeJSON(w, http.StatusOK, map[string]string{"message": "you are now following the user"})
		return
	}

	// rivate profile we create a follow request
	res, err := db.Database.Exec(
		"INSERT INTO follow_user_requests (userToFollow_id, requester_id) VALUES (?, ?)",
		userToFollowID, currentUserID,
	)
	if err != nil {
		log.Println("DB error inserting follow request:", err)
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}

	// ✅ Get the follow request ID
	followRequestID, err := res.LastInsertId()
	if err != nil {
		log.Println("DB error getting follow request id:", err)
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}

	// notification for follow request with request ID
	err = generalfuncs.CreateNotification(models.Notification{
		UserID:          userToFollowID,
		Type:            "follow_request",
		SenderID:        &currentUserID,
		FollowRequestID: &followRequestID,
		SenderName:      &senderName,
		SenderAvatar:    &senderAvatar,
	})
	if err != nil {
		log.Println("DB error inserting follow request notification:", err)
		// non-fatal error, so we don't return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "follow request sent"})
}

func ViewFollowRequestsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	currentUserID := middleware.GetUserId(middleware.GetCookieValue(r))
	if currentUserID == 0 {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	rows, err := db.Database.Query(`
		SELECT fur.id, fur.requester_id, u.username, fur.status, fur.created_at
		FROM follow_user_requests fur
		LEFT JOIN users u ON u.id = fur.requester_id
		WHERE fur.userToFollow_id = ? AND fur.status = 'pending'
		ORDER BY fur.created_at DESC
	`, currentUserID)
	if err != nil {
		log.Println("DB error querying follow requests:", err)
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}
	defer rows.Close()

	var requests []FollowRequest

	for rows.Next() {
		var req FollowRequest
		// NOTE: scanning directly into time.Time because SQLite returns TEXT timestamps
		var createdAt time.Time

		err := rows.Scan(&req.ID, &req.RequesterID, &req.SenderName, &req.Status, &createdAt)
		if err != nil {
			log.Println("DB scan error:", err)
			writeError(w, http.StatusInternalServerError, "database error")
			return
		}

		req.CreatedAt = createdAt
		requests = append(requests, req)
	}

	writeJSON(w, http.StatusOK, requests)
}

// UpdateFollowRequestStatusHandler approves or rejects follow requests
func UpdateFollowRequestStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	currentUserID := middleware.GetUserId(middleware.GetCookieValue(r))
	if currentUserID == 0 {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	requestIDStr := r.URL.Query().Get("follow_request_id")
	requestID, err := strconv.Atoi(requestIDStr)
	if err != nil || requestID == 0 {
		writeError(w, http.StatusBadRequest, "invalid follow_request_id")
		return
	}

	action := r.URL.Query().Get("action")
	if action != "approve" && action != "reject" {
		writeError(w, http.StatusBadRequest, "invalid action")
		return
	}

	tx, err := db.Database.Begin()
	if err != nil {
		log.Println("DB begin tx error:", err)
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}
	defer tx.Rollback()

	var requesterID int
	var currentStatus string
	err = tx.QueryRow(
		"SELECT requester_id, status FROM follow_user_requests WHERE id = ? AND userToFollow_id = ?",
		requestID, currentUserID,
	).Scan(&requesterID, &currentStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "follow request not found")
			return
		}
		log.Println("DB error selecting follow request:", err)
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}

	if currentStatus != "pending" {
		writeError(w, http.StatusConflict, "request already processed")
		return
	}

	newStatus := "rejected"

	var notifyRequester bool     // CHANGED: defer notification until after commit
	var approverFirstName string // CHANGED: sender should be the approver

	if action == "approve" {
		err = tx.QueryRow(
			"SELECT first_name FROM users WHERE id = ?",
			currentUserID, // CHANGED: fetch approver name, not requester
		).Scan(&approverFirstName)
		if err != nil {
			log.Println("DB error fetching approver name:", err)
			writeError(w, http.StatusInternalServerError, "database error")
			return
		}

		newStatus = "approved"
		notifyRequester = true // CHANGED: mark notification for after commit
	}

	_, err = tx.Exec(
		"UPDATE follow_user_requests SET status = ? WHERE id = ?",
		newStatus, requestID,
	)
	if err != nil {
		log.Println("DB error updating request status:", err)
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}

	if newStatus == "approved" {
		_, err := tx.Exec(
			"INSERT OR IGNORE INTO followers (follower_id, followed_id) VALUES (?, ?)",
			requesterID, currentUserID,
		)
		if err != nil {
			log.Println("DB error inserting follower:", err)
			writeError(w, http.StatusInternalServerError, "database error")
			return
		}

		_, _ = tx.Exec(
			"UPDATE follow_user_requests SET status = 'rejected' WHERE userToFollow_id = ? AND requester_id = ? AND id != ? AND status = 'pending'",
			currentUserID, requesterID, requestID,
		)
	}
	// delete the follow request after processing and the notification so it doesnt get refetched
	_, err = tx.Exec(
		"DELETE FROM follow_user_requests WHERE id = ? AND userToFollow_id = ?",
		requestID, currentUserID,
	)
	_, err = tx.Exec(
		"DELETE FROM notifications WHERE follow_request_id = ?",
		requestID,
	)
	if err != nil {
		log.Println("DB error deleting follow request:", err)
	}

	if err := tx.Commit(); err != nil {
		log.Println("DB commit error:", err)
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}

	if notifyRequester { // CHANGED: notification sent AFTER commit
		err = generalfuncs.CreateNotification(models.Notification{
			UserID:     requesterID,
			Type:       "user_accepted_follow",
			SenderID:   &currentUserID,
			SenderName: &approverFirstName,
		})
		if err != nil {
			log.Println("DB error inserting follow request approved notification:", err)
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "follow request " + newStatus})
}

// UnfollowHandler removes a follower relationship
func UnfollowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	currentUserID := middleware.GetUserId(middleware.GetCookieValue(r))
	if currentUserID == 0 {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID == 0 {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if userID == currentUserID {
		writeError(w, http.StatusBadRequest, "cannot unfollow yourself")
		return
	}

	res, err := db.Database.Exec(
		"DELETE FROM followers WHERE follower_id = ? AND followed_id = ?",
		currentUserID, userID,
	)
	if err != nil {
		log.Println("DB error unfollow:", err)
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		writeError(w, http.StatusNotFound, "not following user")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "unfollowed successfully",
	})
}

// CancelFollowRequestHandler allows requester to cancel a pending follow request
func CancelFollowRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	currentUserID := middleware.GetUserId(middleware.GetCookieValue(r))
	if currentUserID == 0 {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID == 0 {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	//im and idiot also delete te notification for the follow request
	_, err = db.Database.Exec(`
		DELETE FROM notifications
		WHERE user_id = ? AND type = 'follow_request' AND sender_id = ?
	`, userID, currentUserID)
	if err != nil {
		log.Println("DB error deleting follow request notification:", err)
	}
	res, err := db.Database.Exec(`
		DELETE FROM follow_user_requests
		WHERE requester_id = ? AND userToFollow_id = ? AND status = 'pending'
	`, currentUserID, userID)
	if err != nil {
		log.Println("DB error cancel follow request:", err)
		writeError(w, http.StatusInternalServerError, "database error")
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		writeError(w, http.StatusNotFound, "no pending request found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "follow request cancelled",
	})
}
