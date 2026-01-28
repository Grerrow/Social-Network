package models

import "time"

// CREATE TABLE notifications (
//     id INTEGER PRIMARY KEY AUTOINCREMENT,
//     user_id INTEGER NOT NULL,
//     follow_request_id INTEGER,
//     type TEXT NOT NULL CHECK(type IN ('follow_request','new_follower','user_accepted_follow','new_comment', 'group_invitation', 'group_join_request_approved', 'event_created','group_join_request','event_invitation')),
//     sender_id INTEGER,
//     sender_name TEXT,
//     sender_avatar TEXT,
//     post_id INTEGER,
//     group_id INTEGER,
//     group_name TEXT,
//     event_id INTEGER,
//     is_read BOOLEAN DEFAULT 0,
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
//     FOREIGN KEY (sender_id) REFERENCES users(id),
//     FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
//     FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
//     FOREIGN KEY (follow_request_id) REFERENCES follow_user_requests(id) ON DELETE CASCADE
// );

type Notification struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	FollowRequestID *int64    `json:"follow_request_id,omitempty"`
	Type            string    `json:"type"`
	SenderID        *int      `json:"sender_id,omitempty"`
	SenderName      *string   `json:"sender_name,omitempty"`
	SenderAvatar    *string   `json:"sender_avatar,omitempty"`
	PostID          *int      `json:"post_id,omitempty"`
	GroupID         *int      `json:"group_id,omitempty"`
	GroupName       *string   `json:"group_name,omitempty"`
	EventID         *int      `json:"event_id,omitempty"`
	EventDate       *int64    `json:"event_date,omitempty"`
	EventTitle      *string   `json:"event_title,omitempty"`
	IsRead          bool      `json:"is_read"`
	CreatedAt       time.Time `json:"created_at"`
}
