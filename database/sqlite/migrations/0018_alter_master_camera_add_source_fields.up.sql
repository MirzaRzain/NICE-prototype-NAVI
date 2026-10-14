ALTER TABLE master_camera ADD COLUMN usecase_config TEXT NOT NULL DEFAULT '[]';
ALTER TABLE master_camera ADD COLUMN channels TEXT;
ALTER TABLE master_camera ADD COLUMN camera_type TEXT;
ALTER TABLE master_camera ADD COLUMN capture_protocol TEXT;
ALTER TABLE master_camera ADD COLUMN rtc_url TEXT;
