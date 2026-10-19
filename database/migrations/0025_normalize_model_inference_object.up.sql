CREATE TABLE IF NOT EXISTS master_inference_object (
    id BIGSERIAL PRIMARY KEY,
    inference_id BIGINT NOT NULL REFERENCES master_inference(id) ON DELETE CASCADE,
    object_id BIGINT NOT NULL REFERENCES master_object(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (inference_id, object_id)
);

INSERT INTO master_inference_object (inference_id, object_id)
SELECT DISTINCT inference_id, object_id
FROM usecase_inference_object
ON CONFLICT (inference_id, object_id) DO NOTHING;

ALTER TABLE usecase_inference_object
    ADD COLUMN IF NOT EXISTS inference_object_id BIGINT REFERENCES master_inference_object(id) ON DELETE CASCADE;

UPDATE usecase_inference_object uio
SET inference_object_id = mio.id
FROM master_inference_object mio
WHERE mio.inference_id = uio.inference_id
  AND mio.object_id = uio.object_id
  AND uio.inference_object_id IS NULL;

CREATE TABLE IF NOT EXISTS model_inference_object (
    id BIGSERIAL PRIMARY KEY,
    model_id BIGINT NOT NULL REFERENCES master_model(id) ON DELETE CASCADE,
    inference_object_id BIGINT NOT NULL REFERENCES master_inference_object(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    UNIQUE (model_id, inference_object_id)
);

INSERT INTO model_inference_object (model_id, inference_object_id)
SELECT DISTINCT old.model_id, uio.inference_object_id
FROM usecase_inference_object_model old
JOIN usecase_inference_object uio ON uio.id = old.usecase_inference_object_id
WHERE uio.inference_object_id IS NOT NULL
ON CONFLICT (model_id, inference_object_id) DO NOTHING;

ALTER TABLE master_model
    DROP COLUMN IF EXISTS usecase_id;
