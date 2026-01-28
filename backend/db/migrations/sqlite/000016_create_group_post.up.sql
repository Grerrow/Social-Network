CREATE TABLE IF NOT EXISTS group_posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    image TEXT,
    created_at INTEGER NOT NULL,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_group_posts_group ON group_posts(group_id);
CREATE INDEX IF NOT EXISTS idx_group_posts_user ON group_posts(user_id);