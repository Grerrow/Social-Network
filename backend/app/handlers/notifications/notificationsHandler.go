package notifications

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"social-network/app/models"

	"social-network/app/middleware"
	"social-network/db"
)

// why ai made my code ugly ??????????
func GetNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserId(middleware.GetCookieValue(r))
	if userID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := db.Database.Query(`
		SELECT
		  id,
		  user_id,
		  type,
		  sender_id,
		  follow_request_id,
		  sender_name,
		  sender_avatar,
		  post_id,
		  group_id,
		  group_name,
		  event_id,
		  is_read,
		  created_at
		FROM notifications
		WHERE user_id = ?
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		http.Error(w, "failed to get notifications", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var notifications []models.Notification

	for rows.Next() {
		var n models.Notification
		var senderID, postID, groupID, eventID sql.NullInt64

		if err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Type,
			&senderID,
			&n.FollowRequestID,
			&n.SenderName,
			&n.SenderAvatar,
			&postID,
			&groupID,
			&n.GroupName,
			&eventID,
			&n.IsRead,
			&n.CreatedAt,
		); err != nil {
			http.Error(w, "failed to scan notification", http.StatusInternalServerError)
			return
		}

		if senderID.Valid {
			senderIDInt := int(senderID.Int64)
			n.SenderID = &senderIDInt
		}
		if postID.Valid {
			postIDInt := int(postID.Int64)
			n.PostID = &postIDInt
		}
		if groupID.Valid {
			groupIDInt := int(groupID.Int64)
			n.GroupID = &groupIDInt
		}
		if eventID.Valid {
			eventIDInt := int(eventID.Int64)
			n.EventID = &eventIDInt
		}

		notifications = append(notifications, n)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func MarkAllNotificationsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserId(middleware.GetCookieValue(r))
	if userID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	_, err := db.Database.Exec(
		`UPDATE notifications SET is_read = 1 WHERE user_id = ?`,
		userID,
	)
	if err != nil {
		http.Error(w, "failed to update notifications", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func RemoveNotificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("ctxUserID").(int)

	notificationID := r.URL.Query().Get("id")
	if notificationID == "" {
		http.Error(w, "missing notification id", http.StatusBadRequest)
		return
	}
	result, err := db.Database.Exec(
		`DELETE FROM notifications WHERE id = ? AND user_id = ?`,
		notificationID,
		userID,
	)
	if err != nil {
		http.Error(w, "failed to delete notification", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "failed to get affected rows", http.StatusInternalServerError)
		return
	}
	if rowsAffected <= 0 {
		http.Error(w, "notification not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
