PRAGMA auto_vacuum = INCREMENTAL;
PRAGMA foreign_keys = ON;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS sessions;

CREATE TABLE IF NOT EXISTS users (
  id text PRIMARY KEY,
  username text NOT NULL DEFAULT '',
  email text NOT NULL DEFAULT '',
  password TEXT NOT NULL,
  salt text NOT NULL,
  confirmation_hash text,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT chk_username_not_empty CHECK (username != ''),
  CONSTRAINT chk_email_not_empty CHECK (email != ''),
  UNIQUE (email),
  UNIQUE (username)
);

CREATE TRIGGER IF NOT EXISTS tg_users_updated_at
  AFTER UPDATE ON users
  FOR EACH ROW
BEGIN
  UPDATE users SET updated_at = CURRENT_TIMESTAMP
WHERE
  id = old.id;
END;

CREATE TABLE IF NOT EXISTS sessions (
  id text PRIMARY KEY NOT NULL,
  user_id text REFERENCES users(id) ON DELETE CASCADE,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  accessed_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS session_users_idx ON sessions(user_id);
