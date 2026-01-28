CREATE TABLE IF NOT EXISTS group_event_responses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    response TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    FOREIGN KEY (event_id) REFERENCES group_events(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(event_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_group_event_responses_event ON group_event_responses(event_id);