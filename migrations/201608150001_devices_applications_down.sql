ALTER TABLE devices_applications
    DROP COLUMN name,
    DROP COLUMN identifier,
    DROP COLUMN short_version,
    DROP COLUMN version,
    DROP COLUMN bundle_size,
    DROP COLUMN dynamic_size,
    DROP COLUMN is_validated;

ALTER TABLE devices_applications ALTER application_uuid DROP DEFAULT;
-- ALTER TABLE devices_applications ADD PRIMARY KEY (application_uuid); // drop?