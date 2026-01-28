package authorization

import (
	"encoding/json"
	"log"
	"net/http"

	"social-network/app/models"
	"social-network/db"

	"golang.org/x/crypto/bcrypt"
)

//=================TO DO==================
// add input validation and sanitization
// implement email verification
// add rate limiting to prevent abuse
// implement CAPTCHA to prevent bot registrations
// add unit tests for registration functionality
//========================================

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var registerData models.RegisterStruct
	if err := json.NewDecoder(r.Body).Decode(&registerData); err != nil {
		log.Printf("Register: Failed to decode body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Register: Received data - email: %s, username: %s, avatar: %s",
		registerData.Email, registerData.Username, registerData.Avatar)

	if registerData.Email == "" || registerData.Password == "" || registerData.FirstName == "" || registerData.LastName == "" || registerData.DateOfBirth == "" {
		log.Printf("Register: Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Register: Failed to hash password: %v", err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Check if email already exists
	var existingEmail string
	err = db.Database.QueryRow("SELECT email FROM users WHERE email = ?", registerData.Email).Scan(&existingEmail)
	if err == nil {
		log.Printf("Register: Email already exists: %s", registerData.Email)
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	// Check if username already exists (only if username is provided)
	if registerData.Username != "" {
		var existingUsername string
		err = db.Database.QueryRow("SELECT username FROM users WHERE username = ?", registerData.Username).Scan(&existingUsername)
		if err == nil {
			log.Printf("Register: Username already exists: %s", registerData.Username)
			http.Error(w, "Username already taken", http.StatusConflict)
			return
		}
	}

	query := `
		INSERT INTO users (
			email, password, first_name, last_name,
			date_of_birth, avatar, username, about_me,
			is_private, created_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`

	result, err := db.Database.Exec(
		query,
		registerData.Email,
		string(hashedPassword),
		registerData.FirstName,
		registerData.LastName,
		registerData.DateOfBirth,
		registerData.Avatar,
		registerData.Username,
		registerData.AboutMe,
		registerData.IsPrivate,
	)
	if err != nil {
		log.Printf("Register: Failed to insert user: %v", err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	userID, _ := result.LastInsertId()
	log.Printf("Register: User created successfully with ID: %d", userID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Registration successful",
		"user": map[string]interface{}{
			"id":       userID,
			"email":    registerData.Email,
			"username": registerData.Username,
		},
	})
}
