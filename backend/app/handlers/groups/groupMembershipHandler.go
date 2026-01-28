package groups

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"social-network/app/generalfuncs"
	"social-network/app/models"
	"social-network/db"
)

type InviteRequest struct {
	GroupID   int `json:"group_id"`
	InviteeID int `json:"user_id"`
}

type RespondToInviteRequest struct {
	GroupID int `json:"group_id"`
}

func InviteUserToGroup(w http.ResponseWriter, r *http.Request) {
	inviterID, ok := r.Context().Value("ctxUserID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req InviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.GroupID <= 0 || req.InviteeID <= 0 {
		http.Error(w, "Invalid group or user ID", http.StatusBadRequest)
		return
	}

	var isMember bool
	if err := db.Database.QueryRow(
		`SELECT EXISTS(
			SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?
		)`,
		req.GroupID, inviterID,
	).Scan(&isMember); err != nil {
		log.Println("[Invite] membership check failed:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if !isMember {
		http.Error(w, "You must be a member to invite users", http.StatusForbidden)
		return
	}

	// check if invitee is already a member
	var alreadyMember bool
	if err := db.Database.QueryRow(
		`SELECT EXISTS(
			SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?
		)`,
		req.GroupID, req.InviteeID,
	).Scan(&alreadyMember); err != nil {
		log.Println("[Invite] member check failed:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if alreadyMember {
		http.Error(w, "User is already a member", http.StatusBadRequest)
		return
	}

	// check if invitation already exists
	var existingInvite bool
	if err := db.Database.QueryRow(
		`SELECT EXISTS(
			SELECT 1 FROM group_invitations
			WHERE group_id = ? AND user_id = ? AND status = 'pending'
		)`,
		req.GroupID, req.InviteeID,
	).Scan(&existingInvite); err != nil {
		log.Println("[Invite] invite check failed:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if existingInvite {
		http.Error(w, "Invitation already sent", http.StatusBadRequest)
		return
	}

	now := time.Now().Unix()

	// create invitation
	if _, err := db.Database.Exec(
		`INSERT INTO group_invitations
		 (group_id, inviter_id, user_id, status, created_at)
		 VALUES (?, ?, ?, 'pending', ?)`,
		req.GroupID, inviterID, req.InviteeID, now,
	); err != nil {
		log.Println("[Invite] insert invitation failed:", err)
		http.Error(w, "Failed to send invitation", http.StatusInternalServerError)
		return
	}

	// query groupname
	var groupName string
	if err := db.Database.QueryRow(
		"SELECT group_name FROM groups WHERE id = ?",
		req.GroupID,
	).Scan(&groupName); err != nil {
		log.Println("[Invite] group lookup failed:", err)
		http.Error(w, "Group not found", http.StatusInternalServerError)
		return
	}

	// query inviter name
	var inviterName string
	if err := db.Database.QueryRow(
		"SELECT username FROM users WHERE id = ?",
		inviterID,
	).Scan(&inviterName); err != nil {
		log.Println("[Invite] inviter lookup failed:", err)
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// Create notification
	// send ws notification
	err := generalfuncs.CreateNotification(models.Notification{
		UserID:     req.InviteeID,
		Type:       "group_invitation",
		SenderID:   &inviterID,
		SenderName: &inviterName,
		GroupID:    &req.GroupID,
		GroupName:  &groupName,
		CreatedAt:  time.Now(),
	})
	if err != nil {
		log.Println("[Invite] notification creation failed:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Invitation sent",
	})
}

func AcceptGroupInvite(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("ctxUserID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// FIXED: Decode into the struct, not into a field
	var req RespondToInviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("[AcceptInvite] decode error:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("[AcceptInvite] userID=%d, groupID=%d", userID, req.GroupID)

	if req.GroupID <= 0 {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	// FIXED: Use user_id column (not invitee_id) to match your INSERT
	var inviteID int
	err := db.Database.QueryRow(
		"SELECT id FROM group_invitations WHERE group_id = ? AND user_id = ? AND status = 'pending'",
		req.GroupID, userID,
	).Scan(&inviteID)

	if err != nil {
		log.Println("[AcceptInvite] invitation not found:", err)
		http.Error(w, "Invitation not found", http.StatusNotFound)
		return
	}

	now := time.Now().Unix()

	// Update invitation status
	_, err = db.Database.Exec(
		"UPDATE group_invitations SET status = 'accepted' WHERE id = ?",
		inviteID,
	)
	if err != nil {
		log.Println("[AcceptInvite] update failed:", err)
	}

	// Add user as member
	_, err = db.Database.Exec(
		"INSERT INTO group_members (group_id, user_id, role, joined_at) VALUES (?, ?, 'member', ?)",
		req.GroupID, userID, now,
	)

	if err != nil {
		log.Println("[AcceptInvite] insert member failed:", err)
		http.Error(w, "Failed to join group", http.StatusInternalServerError)
		return
	}

	// delete notification and invitation
	db.Database.Exec(
		"DELETE FROM notifications WHERE user_id = ? AND type = 'group_invitation' AND group_id = ?",
		userID, req.GroupID,
	)

	db.Database.Exec(
		"DELETE FROM group_invitations WHERE user_id = ? AND group_id = ?",
		userID, req.GroupID,
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Joined group successfully"})
}

func DeclineGroupInvite(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("ctxUserID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// FIXED: Decode into struct, not into int directly
	var req RespondToInviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("[DeclineInvite] decode error:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("[DeclineInvite] userID=%d, groupID=%d", userID, req.GroupID)

	if req.GroupID <= 0 {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	// FIXED: Use user_id column (not invitee_id) to match your INSERT
	result, err := db.Database.Exec(
		"UPDATE group_invitations SET status = 'declined' WHERE group_id = ? AND user_id = ? AND status = 'pending'",
		req.GroupID, userID,
	)

	if err != nil {
		log.Println("[DeclineInvite] update failed:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Invitation not found", http.StatusNotFound)
		return
	}

	// Delete related notifications
	//and delete the invitation itself
	db.Database.Exec(
		"DELETE FROM notifications WHERE user_id = ? AND type = 'group_invitation' AND group_id = ?",
		userID, req.GroupID,
	)

	db.Database.Exec(
		"DELETE FROM group_invitations WHERE user_id = ? AND group_id = ?",
		userID, req.GroupID,
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Invitation declined"})
}

func RequestJoinGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("ctxUserID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	groupID := r.FormValue("group_id")
	log.Printf("[JoinRequest] userID=%d, groupID=%s", userID, groupID)

	if groupID == "" {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	groupIDInt, _ := strconv.Atoi(groupID)

	// Check if already a member
	var isMember bool
	db.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?)",
		groupID, userID,
	).Scan(&isMember)

	if isMember {
		http.Error(w, "Already a member", http.StatusBadRequest)
		return
	}

	// Check if request already exists
	var existingRequest bool
	db.Database.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM group_join_requests WHERE group_id = ? AND user_id = ? AND status = 'pending')",
		groupID, userID,
	).Scan(&existingRequest)

	if existingRequest {
		http.Error(w, "Request already pending", http.StatusBadRequest)
		return
	}

	now := time.Now().Unix()
	// CREATE TABLE IF NOT EXISTS group_join_requests (
	//     id INTEGER PRIMARY KEY AUTOINCREMENT,
	//     group_id INTEGER NOT NULL,
	//     user_id INTEGER NOT NULL,
	//     status TEXT DEFAULT 'pending',
	//     created_at INTEGER NOT NULL,
	//     FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
	//     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	//     UNIQUE(group_id, user_id)
	// );
	_, err := db.Database.Exec(
		"INSERT INTO group_join_requests (group_id, user_id, status, created_at) VALUES (?, ?, 'pending', ?)",
		groupIDInt, userID, now,
	)

	if err != nil {
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		return
	}

	// Get group creator and title
	var creatorID int
	var groupName string
	err = db.Database.QueryRow("SELECT COALESCE(group_name, title), creator_id FROM groups WHERE id = ?", groupID).Scan(&groupName, &creatorID)
	if err != nil {
		http.Error(w, "Failed to get group info", http.StatusInternalServerError)
		return
	}

	//get sender name
	var senderName string
	err = db.Database.QueryRow("SELECT COALESCE(username, first_name) FROM users WHERE id = ?", userID).Scan(&senderName)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	// notification and ws
	err = generalfuncs.CreateNotification(models.Notification{
		UserID:     creatorID,
		Type:       "group_join_request",
		SenderID:   &userID,
		SenderName: &senderName,
		GroupID:    &groupIDInt,
		GroupName:  &groupName,
		CreatedAt:  time.Now(),
	})
	if err != nil {
		log.Println("[JoinRequest] notification creation failed:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Request sent"})
}

func GetJoinRequests(w http.ResponseWriter, r *http.Request) {
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

	// Check if user is the creator
	var creatorID int
	db.Database.QueryRow("SELECT creator_id FROM groups WHERE id = ?", groupID).Scan(&creatorID)

	if creatorID != userID {
		http.Error(w, "Only the creator can view join requests", http.StatusForbidden)
		return
	}

	rows, err := db.Database.Query(`
        SELECT gjr.id, gjr.user_id, gjr.created_at, u.username, u.avatar
        FROM group_join_requests gjr
        JOIN users u ON gjr.user_id = u.id
        WHERE gjr.group_id = ? AND gjr.status = 'pending'
        ORDER BY gjr.created_at DESC
    `, groupID)

	if err != nil {
		http.Error(w, "Failed to fetch requests", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type JoinRequest struct {
		ID        int    `json:"id"`
		UserID    int    `json:"user_id"`
		Username  string `json:"username"`
		Avatar    string `json:"avatar,omitempty"`
		CreatedAt int64  `json:"created_at"`
	}

	requests := []JoinRequest{}
	for rows.Next() {
		var req JoinRequest
		var avatar *string
		err := rows.Scan(&req.ID, &req.UserID, &req.CreatedAt, &req.Username, &avatar)
		if err != nil {
			continue
		}
		if avatar != nil {
			req.Avatar = *avatar
		}
		requests = append(requests, req)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

type ApproveRequestPayload struct {
	GroupID     int `json:"group_id"`
	RequesterID int `json:"user_id"`
}

func ApproveJoinRequest(w http.ResponseWriter, r *http.Request) {
	approverID, ok := r.Context().Value("ctxUserID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req ApproveRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.GroupID <= 0 || req.RequesterID <= 0 {
		http.Error(w, "Group ID and User ID are required", http.StatusBadRequest)
		return
	}

	// Check creator
	var creatorID int
	err := db.Database.QueryRow(
		"SELECT creator_id FROM groups WHERE id = ?",
		req.GroupID,
	).Scan(&creatorID)

	if err == sql.ErrNoRows {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if creatorID != approverID {
		http.Error(w, "Only the creator can approve requests", http.StatusForbidden)
		return
	}

	// Check join request exists
	var requestID int
	var status string
	err = db.Database.QueryRow(
		"SELECT id, status FROM group_join_requests WHERE group_id = ? AND user_id = ?",
		req.GroupID,
		req.RequesterID,
	).Scan(&requestID, &status)

	if err == sql.ErrNoRows {
		http.Error(w, "Join request not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if status != "pending" {
		http.Error(w, "Join request is not pending", http.StatusBadRequest)
		return
	}

	// transaction explanation: begin, rollback on error, commit at end
	//very useful if we want to ensure all db operations succeed or none at all
	//if one fails we rollback to the initial state
	//else we commit the changes
	tx, err := db.Database.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	now := time.Now().Unix()

	// Approve request using the actual request ID
	if _, err := tx.Exec(
		"UPDATE group_join_requests SET status = 'approved' WHERE id = ?",
		requestID,
	); err != nil {
		http.Error(w, "Failed to approve request", http.StatusInternalServerError)
		return
	}

	// Add member
	if _, err := tx.Exec(
		`INSERT INTO group_members (group_id, user_id, role, joined_at)
         VALUES (?, ?, 'member', ?)`,
		req.GroupID,
		req.RequesterID,
		now,
	); err != nil {
		http.Error(w, "Failed to add group member", http.StatusInternalServerError)
		return
	}

	// Commit transaction FIRST before sending notifications
	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// Get group name for notification (after commit)
	//get group name if null get the title instead
	//coalesce returns the first non-null value
	//useful i may say
	var groupName string
	err = db.Database.QueryRow(
		"SELECT COALESCE(group_name, title) FROM groups WHERE id = ?",
		req.GroupID,
	).Scan(&groupName)
	if err != nil {
		http.Error(w, "Failed to get group name", http.StatusInternalServerError)
		log.Println("[ApproveJoinRequest] group name fetch error:", err)
		return
	}

	// Send notification to the requester that their request was approved
	// have to test if using a go routine here is safe, more efficient or in general a good idea
	go func() {
		generalfuncs.CreateNotification(models.Notification{
			UserID: req.RequesterID,
			Type:   "group_join_request_approved", // have to add this type
			//cause i forgot it before
			SenderID:  &approverID,
			GroupID:   &req.GroupID,
			GroupName: &groupName,
			CreatedAt: time.Now(),
		})
	}()

	// always delete the join request notification and  request after approval
	//fooled me once shame on you
	//fooled me twice shame on me
	_, err = db.Database.Exec(
		"DELETE FROM notifications WHERE group_id = ? AND sender_id = ? AND type = 'group_join_request'",
		req.GroupID, req.RequesterID,
	)
	if err != nil {
		http.Error(w, "Failed to delete notification", http.StatusInternalServerError)
		return
	}
	_, err = db.Database.Exec(
		"DELETE FROM group_join_requests WHERE group_id = ? AND user_id = ?",
		req.GroupID, req.RequesterID,
	)
	if err != nil {
		http.Error(w, "Failed to delete group join request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Request approved successfully",
	})
}

func RejectJoinRequest(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("ctxUserID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req ApproveRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.GroupID <= 0 || req.RequesterID <= 0 {
		http.Error(w, "Group ID and Requester ID are required", http.StatusBadRequest)
		return
	}

	// Check if user is the creator
	var creatorID int
	db.Database.QueryRow("SELECT creator_id FROM groups WHERE id = ?", req.GroupID).Scan(&creatorID)

	if creatorID != userID {
		http.Error(w, "Only the creator can reject requests", http.StatusForbidden)
		return
	}

	_, err := db.Database.Exec(
		"DELETE FROM notifications WHERE group_id = ? AND sender_id = ? AND type = 'group_join_request'",
		req.GroupID, req.RequesterID,
	)
	if err != nil {
		http.Error(w, "Failed to delete notification", http.StatusInternalServerError)
		return
	}
	_, err = db.Database.Exec(
		"DELETE FROM group_join_requests WHERE group_id = ? AND user_id = ?",
		req.GroupID, req.RequesterID,
	)
	if err != nil {
		http.Error(w, "Failed to delete group join request", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Request rejected"})
}

func LeaveGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("ctxUserID").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		GroupID int `json:"group_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.GroupID <= 0 {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	// Remove member
	result, err := db.Database.Exec(
		"DELETE FROM group_members WHERE group_id = ? AND user_id = ?",
		req.GroupID, userID,
	)

	if err != nil {
		http.Error(w, "Failed to leave group", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "You are not a member of this group", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Left group successfully"})
}
