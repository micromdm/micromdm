-- +goose Up
CREATE TABLE IF NOT EXISTS cursors (
    value VARCHAR(128) PRIMARY KEY,
		created_at TIMESTAMPTZ DEFAULT '1970-01-01 00:00:00+00'
);


CREATE TABLE IF NOT EXISTS dep_auto_assign (
    profile_uuid VARCHAR(128) PRIMARY KEY,
    filter TEXT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS dep_tokens (
    consumer_key VARCHAR(36) PRIMARY KEY,
    consumer_secret TEXT NULL,
    access_token TEXT NULL,
    access_secret TEXT NULL,
    access_token_expiry TIMESTAMP DEFAULT '1970-01-01 00:00:00'
);

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

CREATE TABLE IF NOT EXISTS push_info (
    udid TEXT PRIMARY KEY,
    token TEXT DEFAULT '',
    push_magic TEXT DEFAULT '',
    mdm_topic TEXT DEFAULT ''
);

CREATE TABLE IF NOT EXISTS server_config (
    config_id INT PRIMARY KEY,
    push_certificate bytea DEFAULT NULL,
    private_key bytea DEFAULT NULL
);

CREATE SEQUENCE IF NOT EXISTS scep_certificates_scep_id_seq
    INCREMENT 1
    MINVALUE 1
    NO MAXVALUE
    START 2;

CREATE TABLE IF NOT EXISTS scep_certificates (
    scep_id integer PRIMARY KEY DEFAULT nextval('scep_certificates_scep_id_seq'),
    cert_name TEXT NULL,
    scep_cert bytea DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS server_config (
    config_id INT PRIMARY KEY,
    push_certificate bytea DEFAULT NULL,
    private_key bytea DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS uuid_cert_auth (
    udid VARCHAR(40) PRIMARY KEY,
    cert_auth bytea DEFAULT NULL
);

-- +goose Down
DROP TABLE IF EXISTS cursors;
DROP TABLE IF EXISTS dep_auto_assign;
DROP TABLE IF EXISTS dep_tokens;
DROP TABLE IF EXISTS devices;
DROP TABLE IF EXISTS push_info;
DROP TABLE IF EXISTS scep_certificates;
DROP TABLE IF EXISTS server_config;
DROP TABLE IF EXISTS uuid_cert_auth;
