DROP INDEX IF EXISTS serial_idx;

-- Need a composite constraint because DEP uses serial and OTA enrollment uses UDID
CREATE UNIQUE INDEX IF NOT EXISTS udid_serial_idx ON devices (udid, serial_number);
