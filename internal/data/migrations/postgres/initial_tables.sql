--DROP TABLE IF EXISTS users;

CREATE TABLE IF NOT EXISTS users (
  id text PRIMARY KEY NOT NULL,
  username text NOT NULL DEFAULT '',
  email text NOT NULL DEFAULT '',
  password bytea NOT NULL,
  salt bytea NOT NULL,
  confirmation_hash text,
  created_at timestamptz DEFAULT (now() at time zone 'utc'),
  updated_at timestamptz DEFAULT (now() at time zone 'utc'),
  CHECK (username != ''),
  CHECK (email != ''),
  UNIQUE (email),
  UNIQUE (username)
);
