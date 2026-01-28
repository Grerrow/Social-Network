CREATE TABLE IF NOT EXISTS group_members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    role VARCHAR(20) DEFAULT 'member',
    joined_at INTEGER NOT NULL,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(group_id, user_id)
);

-- composite primary key (2 primary keys: group_id, user_id) means that we cannot have 2 rows with the same group_id AND user_id
-- ensures each user can only be a member of a group once
CREATE INDEX idx_group_members_group_user
ON group_members(group_id, user_id);

-- find the groups a user belongs to for conversation listing etc
CREATE INDEX idx_user_groups
ON group_members(user_id);