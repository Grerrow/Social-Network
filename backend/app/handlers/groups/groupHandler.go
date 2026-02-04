package groups

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"social-network/app/models"
	"social-network/db"
)

type User struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id,omitempty"`
	Username string `json:"username"`
	Avatar   string `json:"avatar,omitempty"`
}

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("ctxUserID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("CreateGroup: user %d creating group %s", userID, req.Groupname)

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Start a transaction
	tx, err := db.Database.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Insert the group
	result, err := tx.Exec(
		"INSERT INTO groups (group_name, title, description, creator_id) VALUES (?, ?, ?, ?)",
		req.Groupname, req.Title, req.Description, userID,
	)
	if err != nil {
		log.Printf("Failed to create group: %v", err)
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		return
	}

	groupID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed to get group ID: %v", err)
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		return
	}

	// Add creator as member
	_, err = tx.Exec(
		"INSERT INTO group_members (group_id, user_id, role, joined_at) VALUES (?, ?, 'creator', CURRENT_TIMESTAMP)",
		groupID, userID,
	)
	if err != nil {
		log.Printf("Failed to add creator as member: %v", err)
		http.Error(w, "Failed to add creator as member", http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Fetch the created group
	group := models.Group{
		ID:          int(groupID),
		Groupname:   req.Groupname, // Add this line
		Title:       req.Title,
		Description: req.Description,
		CreatorID:   userID,
		MemberCount: 1,
		IsMember:    true,
		IsCreator:   true,
		CreatedAt:   time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(group)
}

func ListGroups(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("ctxUserID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("ListGroups: fetching groups for user %d", userID)
	rows, err := db.Database.Query(`
        SELECT 
            g.id, 
            g.group_name,
            g.title, 
            g.description, 
            g.creator_id, 
            g.created_at,
            (SELECT COUNT(*) FROM group_members WHERE group_id = g.id) as member_count,
            EXISTS(SELECT 1 FROM group_members WHERE group_id = g.id AND user_id = ?) as is_member,
            (g.creator_id = ?) as is_creator,
            EXISTS(SELECT 1 FROM group_invitations WHERE group_id = g.id AND user_id = ? AND status = 'pending') as is_invited,
            EXISTS(SELECT 1 FROM group_join_requests WHERE group_id = g.id AND user_id = ? AND status = 'pending') as has_requested
        FROM groups g
        ORDER BY g.created_at DESC
    `, userID, userID, userID, userID)
	if err != nil {
		log.Printf("Failed to fetch groups: %v", err)
		http.Error(w, "Failed to fetch groups", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	groups := []models.Group{}
	for rows.Next() {
		var g models.Group
		var createdAt time.Time
		err := rows.Scan(
			&g.ID, &g.Groupname, &g.Title, &g.Description, &g.CreatorID, &createdAt,
			&g.MemberCount, &g.IsMember, &g.IsCreator, &g.IsInvited, &g.HasRequested,
		)
		if err != nil {
			log.Printf("Failed to scan group: %v", err)
			continue
		}
		g.CreatedAt = createdAt.Unix()
		groups = append(groups, g)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

func GetGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("ctxUserID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get group ID from query parameter
	groupIDStr := r.URL.Query().Get("id")
	if groupIDStr == "" {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid group ID", http.StatusBadRequest)
		return
	}

	log.Printf("GetGroup: fetching group %d for user %d", groupID, userID)

	var g models.Group
	var createdAt time.Time
	err = db.Database.QueryRow(`
        SELECT 
            g.id, 
            g.group_name,
            g.title, 
            g.description, 
            g.creator_id, 
            g.created_at,
            (SELECT COUNT(*) FROM group_members WHERE group_id = g.id) as member_count,
            EXISTS(SELECT 1 FROM group_members WHERE group_id = g.id AND user_id = ?) as is_member,
            (g.creator_id = ?) as is_creator,
            EXISTS(SELECT 1 FROM group_invitations WHERE group_id = g.id AND user_id = ? AND status = 'pending') as is_invited,
            EXISTS(SELECT 1 FROM group_join_requests WHERE group_id = g.id AND user_id = ? AND status = 'pending') as has_requested
        FROM groups g
        WHERE g.id = ?
    `, userID, userID, userID, userID, groupID).Scan(
		&g.ID, &g.Groupname, &g.Title, &g.Description, &g.CreatorID, &createdAt,
		&g.MemberCount, &g.IsMember, &g.IsCreator, &g.IsInvited, &g.HasRequested,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Group not found", http.StatusNotFound)
			return
		}
		log.Printf("Failed to fetch group: %v", err)
		http.Error(w, "Failed to fetch group", http.StatusInternalServerError)
		return
	}
	g.CreatedAt = createdAt.Unix()

	// Fetch creator info
	var creator models.User
	var avatar sql.NullString
	err = db.Database.QueryRow(
		"SELECT id, username, avatar FROM users WHERE id = ?",
		g.CreatorID,
	).Scan(&creator.ID, &creator.Username, &avatar)
	if err == nil {
		if avatar.Valid {
			avatarStr := avatar.String
			creator.Avatar = &avatarStr
		}
		g.Creator = &creator
	}

	// Fetch members if user is a member
	if g.IsMember {
		memberRows, err := db.Database.Query(`
            SELECT u.id, u.username, COALESCE(u.avatar, '') as avatar
            FROM group_members gm
            JOIN users u ON gm.user_id = u.id
            WHERE gm.group_id = ?
            ORDER BY gm.joined_at ASC
        `, groupID)
		if err == nil {
			defer memberRows.Close()
			for memberRows.Next() {
				var m models.User
				if err := memberRows.Scan(&m.ID, &m.Username, &m.Avatar); err == nil {
					g.Members = append(g.Members, m)
				}
			}
		}

		// Fetch pending requests if user is creator
		if g.IsCreator {
			requestRows, err := db.Database.Query(`
				SELECT u.id, u.username, COALESCE(u.avatar, '') as avatar
				FROM group_join_requests gjr
				JOIN users u ON gjr.user_id = u.id
				WHERE gjr.group_id = ? AND gjr.status = 'pending'
				ORDER BY gjr.created_at DESC
			`, groupID)
			if err == nil {
				defer requestRows.Close()
				for requestRows.Next() {
					var u models.User
					var avatar string
					if err := requestRows.Scan(&u.ID, &u.Username, &avatar); err == nil {
						if avatar != "" {
							u.Avatar = &avatar
						}
						g.PendingRequests = append(g.PendingRequests, u)
					}
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(g)
}

func GetGroupMembers(w http.ResponseWriter, r *http.Request) {
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

	// Verify user is a member
	var isMember bool
	db.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?)",
		groupID, userID,
	).Scan(&isMember)

	if !isMember {
		http.Error(w, "Not a member of this group", http.StatusForbidden)
		return
	}

	rows, err := db.Database.Query(`
        SELECT u.id, u.username, u.avatar, gm.role, gm.joined_at
        FROM group_members gm
        JOIN users u ON gm.user_id = u.id
        WHERE gm.group_id = ?
        ORDER BY gm.joined_at ASC
    `, groupID)

	if err != nil {
		http.Error(w, "Failed to fetch members", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Member struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Avatar   string `json:"avatar,omitempty"`
		Role     string `json:"role"`
		JoinedAt int64  `json:"joined_at"`
	}

	members := []Member{}
	for rows.Next() {
		var m Member
		var avatar sql.NullString
		var joinedAt time.Time
		err := rows.Scan(&m.ID, &m.Username, &avatar, &m.Role, &joinedAt)
		if err != nil {
			continue
		}
		if avatar.Valid {
			m.Avatar = avatar.String
		}
		m.JoinedAt = joinedAt.Unix()
		members = append(members, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//fetch 10 random users

	rows, err := db.Database.Query(`
        SELECT id, username, avatar FROM users
        ORDER BY RANDOM()
        LIMIT 10
    `)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Avatar); err == nil {
			users = append(users, u)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
