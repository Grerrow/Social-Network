CREATE TABLE IF NOT EXISTS group_post_comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT,
    image TEXT,
    created_at INTEGER NOT NULL,
    FOREIGN KEY (post_id) REFERENCES group_posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_group_post_comments_post ON group_post_comments(post_id);