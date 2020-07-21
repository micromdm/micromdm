PRAGMA auto_vacuum = INCREMENTAL;

-- DROP TABLE IF EXISTS users;

CREATE TABLE IF NOT EXISTS users (
  id text PRIMARY KEY,
  username text NOT NULL DEFAULT '',
  email text NOT NULL DEFAULT '',
  password TEXT NOT NULL,
  salt text NOT NULL,
  confirmation_hash text,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  CHECK (username != ''),
  CHECK (email != ''),
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
