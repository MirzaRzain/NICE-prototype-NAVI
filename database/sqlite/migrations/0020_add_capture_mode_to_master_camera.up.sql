ALTER TABLE master_camera
ADD COLUMN capture_mode TEXT CHECK (capture_mode IS NULL OR capture_mode IN ('STREAM', 'SNAPSHOT'));
