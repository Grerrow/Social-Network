CREATE TABLE notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    follow_request_id INTEGER,
    type TEXT NOT NULL CHECK(type IN ('follow_request','new_follower','user_accepted_follow','new_comment', 'group_invitation', 'group_join_request_approved', 'event_created','group_join_request','event_invitation')),
    sender_id INTEGER,
    sender_name TEXT,
    sender_avatar TEXT,
    post_id INTEGER,
    group_id INTEGER,
    group_name TEXT,
    event_id INTEGER,
    event_date TIMESTAMP,
	event_title TEXT,
    is_read BOOLEAN DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES group_events(id) ON DELETE CASCADE,
    FOREIGN KEY (follow_request_id) REFERENCES follow_user_requests(id) ON DELETE CASCADE
);

