ALTER TABLE devices
  ALTER COLUMN dep_profile_assign_time DROP DEFAULT,
  ALTER COLUMN dep_profile_push_time DROP DEFAULT,
  ALTER COLUMN dep_profile_assigned_date DROP DEFAULT,
  ALTER COLUMN last_checkin DROP DEFAULT;
