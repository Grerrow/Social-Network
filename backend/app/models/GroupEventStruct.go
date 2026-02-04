package models

type GroupEvent struct {
	ID            int             `json:"id"`
	GroupID       int             `json:"group_id"`
	CreatorID     int             `json:"creator_id"`
	Title         string          `json:"title"`
	Description   string          `json:"description"`
	EventDate     int64           `json:"event_date"`
	CreatedAt     int64           `json:"created_at"`
	Creator       *User           `json:"creator,omitempty"`
	GoingCount    int             `json:"going_count"`
	NotGoingCount int             `json:"not_going_count"`
	UserResponse  string          `json:"user_response,omitempty"`
	Responses     []EventResponse `json:"responses,omitempty"`
}

type EventResponse struct {
	GroupID   int    `json:"group_id"`
	EventID   int    `json:"event_id"`
	UserID    int    `json:"user_id"`
	Response  string `json:"response"`
	CreatedAt int64  `json:"created_at"`
	User      *User  `json:"user,omitempty"`
}
