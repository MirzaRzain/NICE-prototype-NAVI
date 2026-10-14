CREATE TABLE IF NOT EXISTS usecase_config (
    id BIGSERIAL PRIMARY KEY,
    usecase_id BIGINT NOT NULL REFERENCES master_usecase(id) ON DELETE CASCADE,
    key VARCHAR(255) NOT NULL,
    value TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100),
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(100),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(100),
    UNIQUE (usecase_id, key)
);
