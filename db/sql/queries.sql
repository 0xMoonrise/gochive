-- name: GetArchive :one
SELECT * FROM archive_schema.archive WHERE id = $1 LIMIT 1;

-- name: GetArchiveByName :one
SELECT filename FROM archive_schema.archive WHERE filename=name;

-- name: InsertFile :one
INSERT INTO archive_schema.archive(
	filename,
	editorial,
	file)
VALUES($1, $2, $3)
RETURNING id;

-- name: SaveThumbnail :exec
UPDATE archive_schema.archive
SET
	thumbnail_image=$1
WHERE id=$2;

-- name: GetArchivePage :many
SELECT 
	id, 
	filename,
	editorial,
    favorite
FROM archive_schema.archive 
ORDER BY favorite DESC, id
LIMIT $1
OFFSET $2; 

-- name: GetCountArchive :one
SELECT
	count(id)
FROM archive_schema.archive;

-- name: GetThumbnails :many
SELECT 
	id, 
	thumbnail_image 
FROM archive_schema.archive
ORDER BY id;

-- name: SearchArchive :many
SELECT
    id,
    filename,
    editorial,
    favorite
FROM archive_schema.archive
WHERE filename ILIKE '%' || $1 || '%'
ORDER BY favorite DESC, id
LIMIT $2
OFFSET $3;

-- name: GetCountSearch :one
SELECT
    count(id)
FROM archive_schema.archive
WHERE filename ILIKE '%' || $1 || '%';

-- name: SetFavorite :exec
UPDATE archive_schema.archive
SET favorite=$1
WHERE id = $2;

-- name: SetEditFile :exec
UPDATE 
	archive_schema.archive
SET 
	filename=$1, 
	editorial=$2
WHERE id=$3;

-- name: CreateSchema :exec
CREATE SCHEMA IF NOT EXISTS archive_schema;

-- name: CreateArchiveTable :exec
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

