-- +goose Up
ALTER TABLE archive_schema.archive
ADD COLUMN bookmark FLOAT NOT NULL DEFAULT 0.0;

COMMENT ON COLUMN archive_schema.archive.bookmark IS 'Porcentaje de progreso de lectura (0.0 a 1.0)';

-- +goose Down
ALTER TABLE archive_schema.archive
DROP COLUMN bookmark;
