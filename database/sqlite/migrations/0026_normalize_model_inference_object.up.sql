CREATE TABLE IF NOT EXISTS master_inference_object (id INTEGER PRIMARY KEY AUTOINCREMENT, inference_id INTEGER NOT NULL REFERENCES master_inference(id) ON DELETE CASCADE, object_id INTEGER NOT NULL REFERENCES master_object(id) ON DELETE CASCADE, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP);
CREATE UNIQUE INDEX IF NOT EXISTS uq_sqlite_master_inference_object_pair ON master_inference_object(inference_id, object_id);

INSERT OR IGNORE INTO master_inference_object (inference_id, object_id)
SELECT DISTINCT inference_id, object_id FROM usecase_inference_object;

ALTER TABLE usecase_inference_object ADD COLUMN inference_object_id INTEGER REFERENCES master_inference_object(id) ON DELETE CASCADE;

UPDATE usecase_inference_object
SET inference_object_id = (
    SELECT mio.id
    FROM master_inference_object mio
    WHERE mio.inference_id = usecase_inference_object.inference_id
      AND mio.object_id = usecase_inference_object.object_id
)
WHERE inference_object_id IS NULL;

CREATE TABLE IF NOT EXISTS model_inference_object (id INTEGER PRIMARY KEY AUTOINCREMENT, model_id INTEGER NOT NULL REFERENCES master_model(id) ON DELETE CASCADE, inference_object_id INTEGER NOT NULL REFERENCES master_inference_object(id) ON DELETE CASCADE, created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME, deleted_at DATETIME);
CREATE UNIQUE INDEX IF NOT EXISTS uq_sqlite_model_inference_object_pair ON model_inference_object(model_id, inference_object_id) WHERE deleted_at IS NULL;

ALTER TABLE master_model DROP COLUMN usecase_id;

INSERT OR IGNORE INTO model_inference_object (model_id, inference_object_id)
SELECT DISTINCT old.model_id, uio.inference_object_id
FROM usecase_inference_object_model old
JOIN usecase_inference_object uio ON uio.id = old.usecase_inference_object_id
WHERE uio.inference_object_id IS NOT NULL;
