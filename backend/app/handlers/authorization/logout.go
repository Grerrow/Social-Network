package authorization

import (
	"encoding/json"
	"net/http"

	"social-network/db"
)

// logout handler
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the session cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "No active session", http.StatusUnauthorized)
		return
	}

	// Delete the session from the database
	_, err = db.Database.Exec("DELETE FROM sessions WHERE id = ?", cookie.Value)
	if err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	// Invalidate the cookie
	expiredCookie := http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	}
	http.SetCookie(w, &expiredCookie)

	// Return success response as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Logout successful",
	})
}
