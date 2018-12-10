-- +goose Up

# Required for Timestamp defaults to "zero"
SET sql_mode = '';

CREATE TABLE IF NOT EXISTS devices (
    uuid VARCHAR(40) PRIMARY KEY,
    udid VARCHAR(40) DEFAULT '',
    serial_number VARCHAR(12) DEFAULT '',
    os_version TEXT DEFAULT NULL,
    build_version TEXT DEFAULT NULL,
    product_name TEXT DEFAULT NULL,
    imei TEXT DEFAULT NULL,
    meid TEXT DEFAULT NULL,
    push_magic TEXT DEFAULT NULL,
    awaiting_configuration BOOLEAN DEFAULT false,
    token TEXT DEFAULT NULL,
    unlock_token TEXT DEFAULT NULL,
    enrolled BOOLEAN DEFAULT false,
    description TEXT DEFAULT NULL,
    model TEXT DEFAULT NULL,
    model_name TEXT DEFAULT NULL,
    device_name TEXT DEFAULT NULL,
    color TEXT DEFAULT NULL,
    asset_tag TEXT DEFAULT NULL,
    dep_profile_status TEXT DEFAULT NULL,
    dep_profile_uuid TEXT DEFAULT NULL,
    dep_profile_assign_time TIMESTAMP DEFAULT 0,
    dep_profile_push_time TIMESTAMP DEFAULT 0,
    dep_profile_assigned_date TIMESTAMP DEFAULT 0,
    dep_profile_assigned_by TEXT DEFAULT NULL,
	udid_cert_auth TEXT DEFAULT NULL,
    last_seen TIMESTAMP DEFAULT 0
);


CREATE TABLE IF NOT EXISTS push_info (
    udid VARCHAR(40) PRIMARY KEY,
    token TEXT DEFAULT '',
    push_magic TEXT DEFAULT '',
    mdm_topic TEXT DEFAULT ''
);


-- +goose Down
DROP TABLE IF EXISTS devices;
DROP TABLE IF EXISTS push_info;
