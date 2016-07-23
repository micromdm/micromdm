CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Cant have a unique constraint on name or identifier because it is possible to have multiple versions of the same
-- bundle installed.
CREATE TABLE IF NOT EXISTS applications (
  application_uuid uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  name text NOT NULL,
  identifier text,
  short_version text,
  version text,
  bundle_size bigint,
  dynamic_size bigint,
  is_validated bool,

  install_count integer DEFAULT 1,
  UNIQUE(name, version)
);

CREATE INDEX IF NOT EXISTS idx_application_name ON applications (name);
CREATE INDEX IF NOT EXISTS idx_application_identifier ON applications (identifier);

CREATE TABLE IF NOT EXISTS devices_applications (
  device_uuid uuid REFERENCES devices(device_uuid) ON DELETE CASCADE,
  application_uuid uuid REFERENCES applications(application_uuid) ON DELETE CASCADE
);

