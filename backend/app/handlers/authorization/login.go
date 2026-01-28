package authorization

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"social-network/app/middleware"
	"social-network/app/models"
	"social-network/db"

	"golang.org/x/crypto/bcrypt"
) //=================TO DO==================
// add input validation and sanitization
// implement rate limiting to prevent brute-force attacks
// add account lockout mechanism after multiple failed attempts
// implement multi-factor authentication (MFA)
// add unit tests for login functionality
//========================================

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginReq models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if loginReq.Password == "" || loginReq.Email == "" {
		http.Error(w, "Missing credentials", http.StatusBadRequest)
		return
	}

	userID, err := AuthenticateUser(r.Context(), loginReq)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			http.Error(w, "Wrong email or password", http.StatusUnauthorized)
			return
		}
		log.Printf("auth error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	middleware.CreateSession(userID, w)

	// always return JSON for vue frontend
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"user_id": userID,
	})
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func AuthenticateUser(ctx context.Context, loginReq models.LoginRequest) (int, error) {
	// Query password and id in one go
	var hashedPassword string
	var userID int
	query := "SELECT id, password FROM users WHERE email = ? LIMIT 1"
	err := db.Database.QueryRowContext(ctx, query, loginReq.Email).Scan(&userID, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginReq.Password)) != nil {
		return 0, ErrInvalidCredentials
	}

	return userID, nil
}
