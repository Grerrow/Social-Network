CREATE TABLE IF NOT EXISTS group_invitations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    inviter_id INTEGER NOT NULL,
    status TEXT DEFAULT 'pending',
    created_at INTEGER NOT NULL,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (inviter_id) REFERENCES users(id) ON DELETE CASCADE
);


-- sqlite don't support partial unique indexes directly, so we use a workaround with a filtered unique index
CREATE UNIQUE INDEX IF NOT EXISTS idx_group_invitations_pending ON group_invitations(group_id, user_id) WHERE status='pending';

CREATE INDEX IF NOT EXISTS idx_group_invitations_group ON group_invitations(group_id);
CREATE INDEX IF NOT EXISTS idx_group_invitations_user ON group_invitations(user_id);
CREATE INDEX IF NOT EXISTS idx_group_invitations_status ON group_invitations(status);