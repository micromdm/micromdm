-- Add golang's zero value for time.Time as the default column value for last_checkin which allows us to pass time.Time
-- as a value type as per the docs.
ALTER TABLE devices
  ALTER COLUMN dep_profile_assign_time SET DEFAULT '0001-01-01 00:00:00',
  ALTER COLUMN dep_profile_push_time SET DEFAULT '0001-01-01 00:00:00',
  ALTER COLUMN dep_profile_assigned_date SET DEFAULT '0001-01-01 00:00:00',
  ALTER COLUMN last_checkin SET DEFAULT '0001-01-01 00:00:00'


