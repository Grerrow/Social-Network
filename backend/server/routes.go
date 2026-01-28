package server

import (
	"net/http"
	"os"

	"social-network/app/handlers/authorization"
	"social-network/app/handlers/chat"
	"social-network/app/handlers/comment"
	"social-network/app/handlers/groups"
	"social-network/app/handlers/images"
	"social-network/app/handlers/notifications"
	"social-network/app/handlers/post"
	"social-network/app/handlers/profile"
	"social-network/app/handlers/searchbar"
	"social-network/app/handlers/websocket"
	"social-network/app/middleware"
)

func SetupRoutes() http.Handler {
	mux := http.NewServeMux() // HTTP Request Multiplexer (router)

	// ROUTING ====================================================================================
	// Public authentication endpoints
	mux.HandleFunc("/register", authorization.RegisterHandler)
	mux.HandleFunc("/login", authorization.LoginHandler)
	//de-auth user
	mux.HandleFunc("/logout", middleware.RequireAuth(authorization.LogoutHandler))

	// Auth-required endpoints
	mux.HandleFunc("/me", middleware.MeHandler)

	// Profile
	mux.HandleFunc("/profile", middleware.RequireAuth(profile.ProfileHandler))
	mux.HandleFunc("/profile/privacy", middleware.RequireAuth(profile.UpdateProfilePrivacyHandler))

	// Notifications
	mux.HandleFunc("/notifications", middleware.RequireAuth(notifications.GetNotificationsHandler))
	mux.HandleFunc("/notifications/read-all", middleware.RequireAuth(notifications.MarkAllNotificationsRead))
	mux.HandleFunc("/notifications/remove", middleware.RequireAuth(notifications.RemoveNotificationHandler))

	// follow requests
	mux.HandleFunc("/follow/request", middleware.RequireAuth(profile.FollowRequestHandler))
	mux.HandleFunc("/follow/requests/pending", middleware.RequireAuth(profile.ViewFollowRequestsHandler))
	mux.HandleFunc("/follow/request/status", middleware.RequireAuth(profile.UpdateFollowRequestStatusHandler))
	mux.HandleFunc("/follow/unfollow", middleware.RequireAuth(profile.UnfollowHandler))
	mux.HandleFunc("/follow/request/cancel", middleware.RequireAuth(profile.CancelFollowRequestHandler))

	// Upload images
	fs := http.FileServer(http.Dir("./public/images"))
	mux.Handle("/images/", http.StripPrefix("/images/", fs))
	mux.HandleFunc("/image/upload", images.ImageUploadHandler)

	// Posts
	mux.HandleFunc("/posts", middleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			post.GetFeedPosts(w, r)
		case http.MethodPost:
			post.CreatePost(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// single post retrieval
	mux.HandleFunc("/post", middleware.RequireAuth(post.GetPostHandler))
	// getFollowers endpoint for almost private posts
	mux.HandleFunc("/followers", middleware.RequireAuth(post.GetFollowersHandler))
	// WebSocket endpoint
	mux.HandleFunc("/ws", middleware.RequireAuth(websocket.WebSocketHandler))

	// //private chat and group chat messages
	mux.HandleFunc("/chat/messages", middleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			chat.GetMessages(w, r)
		case http.MethodPost:
			chat.SendMessage(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Group chat messages
	mux.HandleFunc("/chat/group-messages", middleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			chat.GetGroupMessages(w, r)
		case http.MethodPost:
			chat.SendGroupMessage(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// // Chat conversations
	mux.HandleFunc("/chat/conversations", middleware.RequireAuth(chat.GetConversations))
	// mux.HandleFunc("/chat/messages", middleware.RequireAuth(chat.GetMessages))

	// Comments
	mux.HandleFunc("/comments", middleware.RequireAuth(comment.CommentsHandler))

	// search bar
	// Search (users / groups)
	mux.HandleFunc("/search", middleware.RequireAuth(searchbar.SearchBarHandler))

	// Groups
	mux.Handle("/groups", middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			groups.CreateGroup(w, r)
		case http.MethodGet:
			groups.ListGroups(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/groups/details", middleware.RequireAuth(http.HandlerFunc(groups.GetGroup)))
	mux.Handle("/groups/members", middleware.RequireAuth(http.HandlerFunc(groups.GetGroupMembers)))

	// get users for inviting to group
	mux.Handle("/users", middleware.RequireAuth(http.HandlerFunc(groups.GetAllUsers)))

	// Group Invitations
	mux.Handle("/groups/invite", middleware.RequireAuth(http.HandlerFunc(groups.InviteUserToGroup)))
	mux.Handle("/groups/invite/accept", middleware.RequireAuth(http.HandlerFunc(groups.AcceptGroupInvite)))
	mux.Handle("/groups/invite/decline", middleware.RequireAuth(http.HandlerFunc(groups.DeclineGroupInvite)))

	// Group Join Requests
	mux.Handle("/groups/request", middleware.RequireAuth(http.HandlerFunc(groups.RequestJoinGroup)))
	mux.Handle("/groups/requests", middleware.RequireAuth(http.HandlerFunc(groups.GetJoinRequests)))
	mux.Handle("/groups/request/approve", middleware.RequireAuth(http.HandlerFunc(groups.ApproveJoinRequest)))
	mux.Handle("/groups/request/reject", middleware.RequireAuth(http.HandlerFunc(groups.RejectJoinRequest)))

	// Group Posts
	mux.Handle("/groups/posts", middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			groups.CreateGroupPost(w, r)
		case http.MethodGet:
			groups.ListGroupPosts(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Group Comments
	mux.Handle("/groups/comments", middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			groups.CreateGroupComment(w, r)
		case http.MethodGet:
			groups.ListGroupComments(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Group Events
	mux.Handle("/groups/events", middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			groups.CreateGroupEvent(w, r)
		case http.MethodGet:
			groups.ListGroupEvents(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/groups/events/respond", middleware.RequireAuth(http.HandlerFunc(groups.RespondToEvent)))
	//added leave group handler
	mux.Handle("/groups/leave", middleware.RequireAuth(http.HandlerFunc(groups.LeaveGroup)))

	// ============================================================================================

	handler := enableCORS(mux)
	return handler
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CHANGED: Use environment variable with fallback
		allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
		if allowedOrigin == "" {
			allowedOrigin = "http://localhost:80"
		}

		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

/* notes:

- setupRoutes() returns a handler and registers all the route handlers
	this handler is actually a wrapper for the "inside" mux router-handler

- EnableCORS() is the function that actually creates the outside handler-wrapper (which is actually
	what setupRoutes() returns)
	It runs only once in main.go when setupRoutes() is called, not on every request
	(http.HandlerFunc takes an anonymous function as an argument and makes it an http.Handler type function)

- in main.go:
	we create a handler with "router := SetupRoutes()"
	then we pass that handler to http.ListenAndServe as the second argument

	now everytime a request comes in, (frontend/browser -> backend), ListenAndServe() automatically runs the
	handler (with "router.ServeHTTP(w, r)" happening in the background)
	this handler, (which is the CORS-enabling handler), does 2 things:
		1. enables CORS headers
		2. calls a second .ServeHTTP(w, r) but for the inside mux router-handler this time
	this inside mux router-handler matches the URL path to the correct handler function (url:/createPost => post.CreatePost)

*/
