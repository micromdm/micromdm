CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS devices_certificates (
  certificate_uuid uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  device_uuid uuid REFERENCES devices(device_uuid) ON DELETE CASCADE,
  common_name text NOT NULL,
  data BYTEA NOT NULL,
  is_identity BOOL DEFAULT false
)

