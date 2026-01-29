package middleware

import (
	"encoding/json"
	"net/http"

	"social-network/db"
)

// for front end to verify session and get user info (useful for having user data on page load)
func MeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Not authenticated"})
		return
	}

	userID := GetUserId(cookie.Value)
	if userID <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": "Invalid session"})
		return
	}

	// Optionally fetch user details
	var username, email, firstName, lastName, avatar string
	err = db.Database.QueryRow("SELECT username, email, avatar, first_name, last_name FROM users WHERE id = ?", userID).Scan(&username, &email, &avatar, &firstName, &lastName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": map[string]interface{}{
			"id":         userID,
			"username":   username,
			"first_name": firstName,
			"last_name":  lastName,
			"email":      email,
			"avatar":     avatar,
		},
	})
}
