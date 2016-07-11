CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS certificates (
  certificate_uuid uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  common_name text NOT NULL,
  data BYTEA NOT NULL,
  is_identity BOOL DEFAULT false
)

