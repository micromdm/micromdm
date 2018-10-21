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
                awaiting_configuration BOOLEAN DEFAULT FALSE,
                token TEXT DEFAULT '',
                unlock_token TEXT DEFAULT '',
                enrolled BOOLEAN DEFAULT FALSE,
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

CREATE TABLE IF NOT EXISTS push_info (
                udid TEXT PRIMARY KEY,
                token TEXT DEFAULT '',
                push_magic TEXT DEFAULT '',
                mdm_topic TEXT DEFAULT ''
);

CREATE TABLE IF NOT EXISTS wf_profiles (
                id TEXT PRIMARY KEY,
                payload_identifier TEXT NOT NULL CHECK (payload_identifier <> ''),
                payload_uuid TEXT NOT NULL CHECK (payload_uuid <> ''),
                mobileconfig_data bytea,
                created_at TIMESTAMP DEFAULT '1970-01-01 00:00:00',
                updated_at TIMESTAMP DEFAULT '1970-01-01 00:00:00',
                created_by TEXT NOT NULL CHECK (created_by <> ''),
                updated_by TEXT NOT NULL CHECK (updated_by <> ''),
                UNIQUE (payload_identifier,
                        payload_uuid)
);

CREATE TABLE IF NOT EXISTS wf_profile_devices (
                udid TEXT PRIMARY KEY,
                certificate bytea,
                created_at TIMESTAMP DEFAULT '1970-01-01 00:00:00',
                updated_at TIMESTAMP DEFAULT '1970-01-01 00:00:00' -- last seen
);

CREATE TABLE IF NOT EXISTS wf_device_profiles (
                id TEXT PRIMARY KEY,
                command_uuid TEXT DEFAULT '',
                profile_id TEXT REFERENCES wf_profiles (id) ON DELETE CASCADE,
                profile_uuid TEXT NOT NULL CHECK (profile_uuid <> ''),
                device_udid TEXT REFERENCES wf_profile_devices (udid) ON DELETE CASCADE,
                created_by_workflow BOOLEAN DEFAULT TRUE,
                acknowledged BOOLEAN DEFAULT FALSE,
                acknowledged_at TIMESTAMP DEFAULT '1970-01-01 00:00:00',
                created_at TIMESTAMP DEFAULT '1970-01-01 00:00:00',
                updated_at TIMESTAMP DEFAULT '1970-01-01 00:00:00',
                last_seen_on_device TIMESTAMP DEFAULT '1970-01-01 00:00:00',
                UNIQUE (profile_uuid,
                        device_udid)
);

CREATE TABLE IF NOT EXISTS wf_profile_ca (
                id TEXT PRIMARY KEY,
                ca_certificate bytea,
                ca_private_key bytea,
                created_at TIMESTAMP DEFAULT '1970-01-01 00:00:00',
                updated_at TIMESTAMP DEFAULT '1970-01-01 00:00:00'
);

-- +goose Down
DROP TABLE IF EXISTS devices;

DROP TABLE IF EXISTS push_info;

DROP TABLE IF EXISTS wf_profiles;

