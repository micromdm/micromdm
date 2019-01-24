-- +goose Up

# Required for Timestamp defaults to "zero"
SET sql_mode = '';

# platform/dep/sync/mysql
CREATE TABLE IF NOT EXISTS cursors (
    value VARCHAR(128) PRIMARY KEY,
    created_at TIMESTAMP DEFAULT 0
);

# platform/dep/sync/mysql
CREATE TABLE IF NOT EXISTS dep_auto_assign (
    profile_uuid VARCHAR(128) PRIMARY KEY,
    filter TEXT DEFAULT NULL
);

# platform/profie/mysql
CREATE TABLE IF NOT EXISTS profiles (
	profile_id INT(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    identifier TEXT DEFAULT NULL,
    mobileconfig BLOB DEFAULT NULL
);

# platform/scep/mysql
CREATE TABLE IF NOT EXISTS server_config (
	config_id INT PRIMARY KEY,
    push_certificate BLOB DEFAULT NULL,
    private_key BLOB DEFAULT NULL
);

# platform/scep/mysql
CREATE TABLE IF NOT EXISTS scep_certificates (
	scep_id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
	cert_name TEXT NULL,
	scep_cert BLOB DEFAULT NULL
) AUTO_INCREMENT=2;

# platform/queue/mysql
CREATE TABLE IF NOT EXISTS device_commands (
    uuid VARCHAR(40) PRIMARY KEY,
    device_udid VARCHAR(40) NOT NULL,
    payload BLOB DEFAULT NULL,
    created_at TIMESTAMP DEFAULT 0,
    last_sent_at TIMESTAMP DEFAULT 0,
    acknowledged_at TIMESTAMP DEFAULT 0,
    times_sent int(11) DEFAULT 0,
    last_status VARCHAR(32) DEFAULT NULL,
    failure_message BLOB DEFAULT NULL,
    command_order int(11) DEFAULT 0
);

# platform/devices/mysql
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
    last_seen TIMESTAMP DEFAULT 0
);

# platform/devices/mysql
CREATE TABLE IF NOT EXISTS uuid_cert_auth (
    udid VARCHAR(40) PRIMARY KEY,
    cert_auth BLOB DEFAULT NULL
);

# platform/remove/mysql
CREATE TABLE IF NOT EXISTS remove_device (
    udid VARCHAR(40) PRIMARY KEY
);

# platform/config/mysql
CREATE TABLE IF NOT EXISTS server_config (
	config_id INT PRIMARY KEY,
    push_certificate BLOB DEFAULT NULL,
    private_key BLOB DEFAULT NULL
);

# platform/config/mysql
CREATE TABLE IF NOT EXISTS dep_tokens (
	consumer_key VARCHAR(36) PRIMARY KEY,
	consumer_secret TEXT NULL,
	access_token TEXT NULL,
	access_secret TEXT NULL,
    access_token_expiry TIMESTAMP DEFAULT 0
);

# platform/apns/mysql
CREATE TABLE IF NOT EXISTS push_info (
    udid VARCHAR(40) PRIMARY KEY,
    token TEXT DEFAULT '',
    push_magic TEXT DEFAULT '',
    mdm_topic TEXT DEFAULT ''
);


-- +goose Down
DROP TABLE IF EXISTS cursors;
DROP TABLE IF EXISTS dep_auto_assign;
DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS server_config;
DROP TABLE IF EXISTS scep_certificates;
DROP TABLE IF EXISTS device_commands;
DROP TABLE IF EXISTS devices;
DROP TABLE IF EXISTS uuid_cert_auth;
DROP TABLE IF EXISTS remove_device;
DROP TABLE IF EXISTS server_config;
DROP TABLE IF EXISTS dep_tokens;
DROP TABLE IF EXISTS push_info;
