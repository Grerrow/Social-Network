CREATE TABLE posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT,
    image TEXT,
    privacy TEXT NOT NULL CHECK(privacy IN ('public', 'almost_private', 'private')) DEFAULT 'public',
    user_id INTEGER NOT NULL,
    group_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);