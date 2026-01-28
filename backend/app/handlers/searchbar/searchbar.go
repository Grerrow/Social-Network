package searchbar

import (
	"encoding/json"
	"net/http"
	"strings"

	"social-network/app/models"
	"social-network/db"
)

// =================TO DO==================
// implement persistent storage (database) for users
// optimize search algorithm for large datasets
// add pagination support for search results
// implement fuzzy search or autocomplete features
// add unit tests for search functionality
// ========================================
// CREATE TABLE users (
//
//	id INTEGER PRIMARY KEY AUTOINCREMENT,
//	email TEXT NOT NULL UNIQUE,
//	password TEXT NOT NULL,
//	first_name TEXT NOT NULL,
//	last_name TEXT NOT NULL,
//	date_of_birth DATE NOT NULL,
//	avatar TEXT,
//	username TEXT,
//	about_me TEXT,
//	is_private BOOLEAN NOT NULL DEFAULT 1,
//	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
//
// );
var users = []models.User{}

// searchbar handler to return users for js to populate the search bar
func SearchBarHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		json.NewEncoder(w).Encode(map[string][]any{
			"users":  {},
			"groups": {},
		})
		return
	}

	query = "%" + query + "%"

	rows, err := db.Database.Query(`
		SELECT id, username, first_name, last_name
		FROM users
		WHERE username LIKE ? OR first_name LIKE ? OR last_name LIKE ?
		LIMIT 10
	`, query, query, query)

	if err != nil {
		http.Error(w, "search failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []map[string]any

	for rows.Next() {
		var id int
		var username, firstName, lastName string

		if err := rows.Scan(&id, &username, &firstName, &lastName); err != nil {
			continue
		}

		users = append(users, map[string]any{
			"id":       id,
			"username": username,
			"name":     firstName + " " + lastName,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"users":  users,
		"groups": []any{}, // placeholder
	})
}

// search registered users by username or full name
func SearchUsers(query string) []models.User {
	var results []models.User
	for _, user := range users {
		fullName := user.FirstName + " " + user.LastName
		username := ""
		if user.Username != nil {
			username = *user.Username
		}
		if containsIgnoreCase(username, query) || containsIgnoreCase(fullName, query) {
			results = append(results, user)
		}
	}
	return results
}

// helper function to check if str contains substr (case insensitive)
func containsIgnoreCase(str, substr string) bool {
	strLower := strings.ToLower(str)
	substrLower := strings.ToLower(substr)
	return strings.Contains(strLower, substrLower)
}
