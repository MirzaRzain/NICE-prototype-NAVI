CREATE TABLE IF NOT EXISTS usecase_inference_object_model (
    id BIGSERIAL PRIMARY KEY,
    model_id BIGINT NOT NULL REFERENCES master_model(id) ON DELETE CASCADE,
    usecase_inference_object_id BIGINT NOT NULL REFERENCES usecase_inference_object(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    UNIQUE (model_id)
);
