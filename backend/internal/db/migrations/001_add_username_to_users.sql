-- Up migration
CREATE TABLE new_users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    first_name TEXT,
    last_name TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO new_users (id, email, username, first_name, last_name, created_at, updated_at)
SELECT id, email, email, '', '', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
FROM users;

DROP TABLE users;
ALTER TABLE new_users RENAME TO users;

-- Down migration
CREATE TABLE old_users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL
);

INSERT INTO old_users (id, email)
SELECT id, email
FROM users;

DROP TABLE users;
ALTER TABLE old_users RENAME TO users;