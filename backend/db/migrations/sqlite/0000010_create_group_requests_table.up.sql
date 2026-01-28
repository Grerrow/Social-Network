CREATE TABLE IF NOT EXISTS group_join_requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    status TEXT DEFAULT 'pending',
    created_at INTEGER NOT NULL,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


-- UNIQUE(group_id, user_id, status='pending') cant do that so...
CREATE UNIQUE INDEX IF NOT EXISTS idx_group_join_requests_pending ON group_join_requests(group_id, user_id) WHERE status='pending';

CREATE INDEX IF NOT EXISTS idx_group_join_requests_group ON group_join_requests(group_id);
CREATE INDEX IF NOT EXISTS idx_group_join_requests_user ON group_join_requests(user_id);
CREATE INDEX IF NOT EXISTS idx_group_join_requests_status ON group_join_requests(status);