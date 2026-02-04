package groups

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"social-network/app/handlers/websocket"
	"social-network/app/models"
	"social-network/db"
)

// description	"hbuhbnyh"
// event_date	-27103790578
// group_id	6
// title	"yhuuyh"
type CreateEventRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	EventDate   int64  `json:"event_date"`
	GroupID     int    `json:"group_id"`
}

func CreateGroupEvent(w http.ResponseWriter, r *http.Request) {
	// --- Auth ---
	userID, ok := r.Context().Value("ctxUserID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// --- Decode request ---
	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.GroupID <= 0 {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if req.EventDate <= time.Now().Unix() {
		http.Error(w, "Event date must be in the future", http.StatusBadRequest)
		return
	}

	// --- Check membership ---
	var isMember bool
	if err := db.Database.QueryRow(
		`SELECT EXISTS (
			SELECT 1 FROM group_members 
			WHERE group_id = ? AND user_id = ?
		)`,
		req.GroupID, userID,
	).Scan(&isMember); err != nil {
		http.Error(w, "Membership check failed", http.StatusInternalServerError)
		return
	}

	if !isMember {
		http.Error(w, "Must be a group member to create events", http.StatusForbidden)
		return
	}

	// --- Start transaction ---
	tx, err := db.Database.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	now := time.Now().Unix()

	// --- Insert event ---
	result, err := tx.Exec(
		`INSERT INTO group_events 
		 (group_id, creator_id, title, description, event_date, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		req.GroupID, userID, req.Title, req.Description, req.EventDate, now,
	)
	if err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	eventID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to get event ID", http.StatusInternalServerError)
		return
	}

	// --- Fetch group name ---
	var groupName string
	if err := tx.QueryRow(
		`SELECT group_name FROM groups WHERE id = ?`,
		req.GroupID,
	).Scan(&groupName); err != nil {
		http.Error(w, "Failed to fetch group info", http.StatusInternalServerError)
		return
	}

	// --- Fetch creator info ---
	var creator models.User
	if err := tx.QueryRow(
		`SELECT username, avatar FROM users WHERE id = ?`,
		userID,
	).Scan(&creator.Username, &creator.Avatar); err != nil {
		http.Error(w, "Failed to fetch creator info", http.StatusInternalServerError)
		return
	}

	// --- Get group members ---
	rows, err := tx.Query(
		`SELECT user_id FROM group_members WHERE group_id = ?`,
		req.GroupID,
	)
	if err != nil {
		http.Error(w, "Failed to fetch group members", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var memberIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err == nil {
			memberIDs = append(memberIDs, id)
		}
	}

	// --- Create notifications ---
	for _, memberID := range memberIDs {
		_, err := tx.Exec(
			`INSERT INTO notifications (
				user_id,
				type,
				sender_id,
				sender_name,
				group_id,
				group_name,
				event_id,
				event_title,
				event_date,
				created_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			memberID,
			"event_created",
			userID,
			creator.Username,
			req.GroupID,
			groupName,
			eventID,
			req.Title,
			req.EventDate,
			now,
		)
		if err != nil {
			http.Error(w, "Failed to create notifications", http.StatusInternalServerError)
			return
		}
	}

	// --- Commit transaction ---
	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// ...existing code...

	// --- WebSocket (best effort) ---
	event := models.GroupEvent{
		ID:          int(eventID),
		Title:       req.Title,
		Description: req.Description,
		EventDate:   req.EventDate,
		CreatedAt:   now,
		Creator:     &creator,
	}

	// FIXED: Send complete notification data
	for _, memberID := range memberIDs {
		websocket.SendToUser(strconv.Itoa(memberID), websocket.WebSocketMessage{
			Type: "event_created",
			Data: map[string]interface{}{
				"id":            0, // Frontend will handle
				"type":          "event_created",
				"sender_id":     userID,
				"sender_name":   creator.Username,
				"sender_avatar": creator.Avatar,
				"group_id":      req.GroupID,
				"group_name":    groupName,
				"event_id":      int(eventID),
				"event_title":   req.Title,
				"event_date":    req.EventDate,
				"is_read":       false,
				"created_at":    now,
			},
		})
	}

	// --- Response ---
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

func ListGroupEvents(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("ctxUserID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID := r.URL.Query().Get("id")

	if groupID == "" {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	// Check if user is a member
	var isMember bool
	db.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?)",
		groupID, userID,
	).Scan(&isMember)

	if !isMember {
		http.Error(w, "Must be a member to view events", http.StatusForbidden)
		return
	}

	rows, err := db.Database.Query(`
        SELECT ge.id, ge.group_id, ge.creator_id, ge.title, ge.description, ge.event_date, ge.created_at,
               u.username, u.avatar,
               (SELECT COUNT(*) FROM group_event_responses WHERE event_id = ge.id AND response = 'going') as going_count,
               (SELECT COUNT(*) FROM group_event_responses WHERE event_id = ge.id AND response = 'not_going') as not_going_count,
               (SELECT response FROM group_event_responses WHERE event_id = ge.id AND user_id = ?) as user_response
        FROM group_events ge
        JOIN users u ON ge.creator_id = u.id
        WHERE ge.group_id = ?
        ORDER BY ge.event_date ASC
    `, userID, groupID)

	if err != nil {
		http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	events := []models.GroupEvent{}
	for rows.Next() {
		var e models.GroupEvent
		var creator models.User
		var avatar, userResponse *string

		err := rows.Scan(
			&e.ID, &e.GroupID, &e.CreatorID, &e.Title, &e.Description, &e.EventDate, &e.CreatedAt,
			&creator.Username, &avatar,
			&e.GoingCount, &e.NotGoingCount, &userResponse,
		)
		if err != nil {
			continue
		}

		creator.ID = e.CreatorID
		if avatar != nil {
			creator.Avatar = avatar
		}
		if userResponse != nil {
			e.UserResponse = *userResponse
		}
		e.Creator = &creator

		// Fetch responses for this event
		e.Responses = fetchEventResponses(e.ID)

		events = append(events, e)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func fetchEventResponses(eventID int) []models.EventResponse {
	rows, err := db.Database.Query(`
        SELECT ger.id, ger.user_id, ger.response, ger.created_at,
               u.username, u.avatar
        FROM group_event_responses ger
        JOIN users u ON ger.user_id = u.id
        WHERE ger.event_id = ?
        ORDER BY ger.created_at ASC
    `, eventID)

	if err != nil {
		return []models.EventResponse{}
	}
	defer rows.Close()

	responses := []models.EventResponse{}
	for rows.Next() {
		var r models.EventResponse
		var user models.User
		var avatar *string

		err := rows.Scan(
			&r.GroupID, &r.UserID, &r.Response, &r.CreatedAt,
			&user.Username, &avatar,
		)
		if err != nil {
			continue
		}

		user.ID = r.UserID
		if avatar != nil {
			user.Avatar = avatar
		}
		r.User = &user
		responses = append(responses, r)
	}

	return responses
}

func RespondToEvent(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("ctxUserID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.EventResponse
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Response == "" {
		http.Error(w, "Response is required", http.StatusBadRequest)
		return
	}
	if req.GroupID == 0 || req.GroupID < 0 {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}
	if req.EventID == 0 || req.EventID < 0 {
		http.Error(w, "Event ID is required", http.StatusBadRequest)
		return
	}

	// Check if user is a member
	var isMember bool
	db.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?)",
		req.GroupID, userID,
	).Scan(&isMember)

	if !isMember {
		http.Error(w, "Must be a member to respond to events", http.StatusForbidden)
		return
	}

	// Check if event exists and belongs to this group
	var eventExists bool
	db.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM group_events WHERE id = ? AND group_id = ?)",
		req.EventID, req.GroupID,
	).Scan(&eventExists)

	if !eventExists {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	if req.Response != "going" && req.Response != "not_going" {
		http.Error(w, "Response must be 'going' or 'not_going'", http.StatusBadRequest)
		return
	}

	now := time.Now().Unix()

	// Upsert response
	_, err := db.Database.Exec(`
        INSERT INTO group_event_responses (event_id, user_id, response, created_at) 
        VALUES (?, ?, ?, ?)
        ON CONFLICT(event_id, user_id) DO UPDATE SET response = ?, created_at = ?
    `, req.EventID, userID, req.Response, now, req.Response, now)

	if err != nil {
		http.Error(w, "Failed to save response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Response saved", "response": req.Response})
}
