DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;

CREATE TABLE IF NOT EXISTS users (
  id text PRIMARY KEY NOT NULL,
  username text NOT NULL DEFAULT '',
  email text NOT NULL DEFAULT '',
  password bytea NOT NULL,
  salt bytea NOT NULL,
  confirmation_hash text,
  created_at timestamptz DEFAULT (now() at time zone 'utc'),
  updated_at timestamptz DEFAULT (now() at time zone 'utc'),
  CONSTRAINT chk_username_not_empty CHECK (username != ''),
  CONSTRAINT chk_email_not_empty CHECK (email != ''),
  UNIQUE (email),
  UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS sessions (
  id text PRIMARY KEY NOT NULL,
  user_id text REFERENCES users(id) ON DELETE CASCADE,
  created_at timestamptz DEFAULT (now() at time zone 'utc'),
  accessed_at timestamptz DEFAULT (now() at time zone 'utc')
);
