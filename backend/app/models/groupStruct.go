package models

type CreateGroupRequest struct {
	Groupname   string `json:"group_name"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Group struct {
	ID              int    `json:"id"`
	Groupname       string `json:"group_name"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	CreatorID       int    `json:"creator_id"`
	CreatedAt       int64  `json:"created_at"`
	MemberCount     int    `json:"member_count"`
	IsMember        bool   `json:"is_member"`
	IsCreator       bool   `json:"is_creator"`
	IsInvited       bool   `json:"is_invited"`
	HasRequested    bool   `json:"has_requested"`
	Creator         *User  `json:"creator,omitempty"`
	Members         []User `json:"members,omitempty"`
	PendingRequests []User `json:"pending_requests,omitempty"`
}
