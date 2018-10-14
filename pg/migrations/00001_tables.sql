-- +goose Up
CREATE TABLE IF NOT EXISTS devices (
    uuid TEXT PRIMARY KEY,
    udid TEXT DEFAULT '',
    serial_number TEXT DEFAULT '',
    os_version TEXT DEFAULT '',
    build_version TEXT DEFAULT '',
    product_name TEXT DEFAULT '',
    imei TEXT DEFAULT '',
    meid TEXT DEFAULT '',
    push_magic TEXT DEFAULT '',
    awaiting_configuration BOOLEAN DEFAULT false,
    token TEXT DEFAULT '',
    unlock_token TEXT DEFAULT '',
    enrolled BOOLEAN DEFAULT false,
    description TEXT DEFAULT '',
    model TEXT DEFAULT '',
    model_name TEXT DEFAULT '',
    device_name TEXT DEFAULT '',
    color TEXT DEFAULT '',
    asset_tag TEXT DEFAULT '',
    dep_profile_status TEXT DEFAULT '',
    dep_profile_uuid TEXT DEFAULT '',
    dep_profile_assign_time TIMESTAMP DEFAULT '1970-01-01 00:00:00',
    dep_profile_push_time TIMESTAMP DEFAULT '1970-01-01 00:00:00',
    dep_profile_assigned_date TIMESTAMP DEFAULT '1970-01-01 00:00:00',
    dep_profile_assigned_by TEXT DEFAULT '',
    last_seen TIMESTAMP DEFAULT '1970-01-01 00:00:00'
);

-- +goose Down
DROP TABLE IF EXISTS devices;
