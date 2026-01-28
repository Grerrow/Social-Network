package middleware

import (
	"context"
	"encoding/json"
	"net/http"
)

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		userID := GetUserId(cookie.Value)

		if err != nil || userID <= 0 {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
			return
		}

		// we add userID to the request context
		// handlers that are wrapped with this middleware in routes
		// can access the userID from the context
		ctx := context.WithValue(r.Context(), "ctxUserID", userID)

		// Call next handler with new context
		next(w, r.WithContext(ctx))
	}
}
