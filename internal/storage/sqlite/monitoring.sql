CREATE TABLE IF NOT EXISTS monitoring_sensors(
  id TEXT PRIMARY KEY,
  cave_code TEXT NOT NULL REFERENCES zones(code),
  label TEXT NOT NULL,
  installed_at TEXT NOT NULL,
  status TEXT NOT NULL,
  version INTEGER NOT NULL DEFAULT 1,
  UNIQUE(cave_code, label)
);
CREATE TABLE IF NOT EXISTS environment_readings(
  id TEXT PRIMARY KEY,
  sensor_id TEXT NOT NULL REFERENCES monitoring_sensors(id),
  wall_temperature REAL NOT NULL,
  relative_humidity REAL NOT NULL,
  outside_temperature REAL NOT NULL,
  outside_humidity REAL NOT NULL,
  dew_point REAL NOT NULL,
  state TEXT NOT NULL,
  recorded_at TEXT NOT NULL,
  created_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS treatment_devices(
  id TEXT PRIMARY KEY,
  cave_code TEXT NOT NULL REFERENCES zones(code),
  serial_number TEXT NOT NULL UNIQUE,
  mode TEXT NOT NULL,
  status TEXT NOT NULL,
  version INTEGER NOT NULL DEFAULT 1,
  last_seen_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS treatment_runs(
  id TEXT PRIMARY KEY,
  device_id TEXT NOT NULL REFERENCES treatment_devices(id),
  reading_id TEXT NOT NULL REFERENCES environment_readings(id),
  action TEXT NOT NULL,
  state TEXT NOT NULL,
  started_at TEXT NOT NULL,
  finished_at TEXT,
  UNIQUE(device_id, reading_id)
);
CREATE INDEX IF NOT EXISTS idx_readings_sensor_time ON environment_readings(sensor_id, recorded_at);
CREATE INDEX IF NOT EXISTS idx_treatment_device_status ON treatment_devices(cave_code, status);
