-- +goose Up
ALTER TABLE archive_schema.archive
DROP COLUMN IF EXISTS file,
DROP COLUMN IF EXISTS thumbnail_image;

-- +goose Down
ALTER TABLE archive_schema.archive
ADD COLUMN IF NOT EXISTS file BYTEA,
ADD COLUMN IF NOT EXISTS thumbnail_image BYTEA;
