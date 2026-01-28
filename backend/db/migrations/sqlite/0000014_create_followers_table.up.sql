CREATE TABLE follow_user_requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    userToFollow_id INTEGER NOT NULL,
    requester_id INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending', -- 'pending', 'approved', 'rejected'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- prevent duplicate follow requests
    UNIQUE (userToFollow_id, requester_id),

    FOREIGN KEY (userToFollow_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (requester_id) REFERENCES users(id)
);

--followers table to track who follows whom
CREATE TABLE followers (
    follower_id INTEGER NOT NULL,
    followed_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- composite primary key (2 primary keys: follower_id, followed_id) ensures each user can only follow another user once
    PRIMARY KEY (follower_id, followed_id),

    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (followed_id) REFERENCES users(id) ON DELETE CASCADE
);