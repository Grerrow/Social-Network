package middleware

import (
	"log"
	"net/http"
	"time"

	"social-network/db"

	"github.com/google/uuid"
)

// create a session
func CreateSession(userID int, w http.ResponseWriter) {
	// deleting any existing session:
	_, err := db.Database.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, "Failed to remove old session", http.StatusInternalServerError)
		return
	}

	// generating a new Session ID and storing it in DB:
	sessionID := uuid.New().String()
	expirationTime := time.Now().Add(24 * time.Hour) // 1-day session

	_, err = db.Database.Exec("INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)",
		sessionID, userID, expirationTime)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// creating a cookie
	cookie := http.Cookie{
		Name:     "session_token",
		Value:    sessionID,
		Expires:  expirationTime,
		HttpOnly: true,                    // we use this to prevent JavaScript access
		SameSite: http.SameSiteStrictMode, // to prevent CSRF attacks
		Path:     "/",                     // this way the cookie is available across the entire site
	}

	// setting the cookie in the user's browser
	http.SetCookie(w, &cookie)
}

// get cookie value (session id) from request
func GetCookieValue(r *http.Request) string {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Printf("Error retrieving session_token cookie: %v", err)
		return ""
	}
	return cookie.Value
}

// check if user exists in DB using his cookie.value (session ID)
// THIS FUNCTION IS NOT NECESSARY WHEN YOU HAVE AUTH.MIDDLEWARE
// too late to change it now tho
func GetUserId(cookieValue string) int { // cookieValue = session ID
	var userId int
	err := db.Database.QueryRow("SELECT user_id FROM sessions WHERE id = ?", cookieValue).Scan(&userId)
	if err != nil {
		log.Printf("GetUserId error: %v", err)
		return -1 // returning -1 is clearer for "not found" result
	}
	return userId
}

func GetUsername(userID int) string {
	var username string
	err := db.Database.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
	if err != nil {
		log.Printf("GetUsername error: %v", err)
		return "" // returning empty string if not found
	}
	return username
}
