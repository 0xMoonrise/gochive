-- +goose Up
CREATE SCHEMA IF NOT EXISTS archive_schema;

CREATE TABLE IF NOT EXISTS archive_schema.archive (
    id SERIAL PRIMARY KEY,
    filename TEXT NOT NULL,
    editorial TEXT NOT NULL,
    cover_page INTEGER NOT NULL DEFAULT 1,
    file BYTEA NOT NULL,
    favorite BOOLEAN NOT NULL DEFAULT FALSE,
    thumbnail_image BYTEA,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
    CONSTRAINT check_cover_page_positive CHECK (cover_page >= 1)
);

CREATE INDEX IF NOT EXISTS idx_archive_id ON archive_schema.archive (id);

-- +goose Down
DROP SCHEMA IF EXISTS archive_schema CASCADE;
DROP INDEX IF EXISTS archive_schema.idx_archive_id;
DROP TABLE IF EXISTS archive_schema.archive;
