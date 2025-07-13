-- +goose Up
CREATE SCHEMA IF NOT EXISTS archive_schema;

CREATE TABLE IF NOT EXISTS archive_schema.archive (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    filename TEXT NOT NULL,
    editorial TEXT NOT NULL,
    cover_page INTEGER NOT NULL DEFAULT 1,
    file BYTEA NOT NULL,
    favorite BOOLEAN NOT NULL DEFAULT FALSE,
    thumbnail_image BYTEA,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
    CONSTRAINT check_cover_page_positive CHECK (cover_page >= 1)
);
-- +goose Down
-- DROP SCHEMA archive_schema CASCADE;
-- DROP TABLE IF EXISTS archive_schema.archive;
