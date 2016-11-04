ALTER TABLE devices_applications DROP CONSTRAINT IF EXISTS devices_applications_application_uuid_fkey;
ALTER TABLE devices_applications ADD PRIMARY KEY (application_uuid);
ALTER TABLE devices_applications ALTER application_uuid SET DEFAULT uuid_generate_v4();

ALTER TABLE devices_applications
    ADD COLUMN name TEXT,
    ADD COLUMN identifier TEXT,
    ADD COLUMN short_version TEXT,
    ADD COLUMN version TEXT,
    ADD COLUMN bundle_size BIGINT,
    ADD COLUMN dynamic_size BIGINT,
    ADD COLUMN is_validated BOOLEAN;
