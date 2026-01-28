package profile

import (
	"encoding/json"
	"net/http"
	"social-network/app/middleware"
	"social-network/db"
)

//change profile privacy /profile/privacy
func UpdateProfilePrivacyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentUserID := middleware.GetUserId(middleware.GetCookieValue(r))
	if currentUserID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var reqBody struct {
		IsPrivate bool `json:"is_private"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err := db.Database.Exec(
		"UPDATE users SET is_private = ? WHERE id = ?",
		reqBody.IsPrivate,
		currentUserID,
	)
	if err != nil {
		http.Error(w, "Failed to update profile privacy", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"is_private": reqBody.IsPrivate,
	})
}

//might add more profile update handlers here in future