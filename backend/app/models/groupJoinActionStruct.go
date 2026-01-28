package models

import "time"

type GroupInvitation struct {
	ID        int       `json:"id"`
	GroupID   int       `json:"group_id"`
	InviterID int       `json:"inviter_id"`
	InviteeID int       `json:"invitee_id"`
	Status    string    `json:"status"`     // 'pending', 'accepted', 'rejected'
	CreatedAt time.Time `json:"created_at"`
}

type GroupJoinRequest struct {
	ID          int       `json:"id"`
	GroupID     int       `json:"group_id"`
	RequesterID int       `json:"requester_id"`
	Status      string    `json:"status"`       // 'pending', 'approved', 'rejected'
	CreatedAt   time.Time `json:"created_at"`
}
