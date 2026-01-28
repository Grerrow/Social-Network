CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    date_of_birth DATE NOT NULL,
    avatar TEXT,
    username TEXT,
    about_me TEXT,
    is_private BOOLEAN NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_first_name ON users(first_name);
CREATE INDEX idx_users_last_name ON users(last_name);
CREATE INDEX idx_users_user_id ON users(id);

-- rows, err := db.Database.Query(`
-- 		SELECT id, username, first_name, last_name
-- 		FROM users
-- 		WHERE username LIKE ? OR first_name LIKE ? OR last_name LIKE ?
-- 		LIMIT 10
-- 	`, query, query, query)
