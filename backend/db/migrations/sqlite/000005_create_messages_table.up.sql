CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER NOT NULL,
    receiver_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_read BOOLEAN DEFAULT 0,
    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY (receiver_id) REFERENCES users(id)
);

-- For inbox / unread counts
CREATE INDEX idx_messages_receiver_is_read
ON messages(receiver_id, is_read);

-- For fetching conversation history
CREATE INDEX idx_messages_sender_receiver
ON messages(sender_id, receiver_id);
CREATE INDEX idx_messages_receiver_sender
ON messages(receiver_id, sender_id);

-- For ordering messages by time
CREATE INDEX idx_messages_receiver_created
ON messages(receiver_id, created_at);

-- mark messages as read inxexx
CREATE INDEX idx_messages_receiver_is_read_created
ON messages(sender_id, receiver_id, is_read);


