-- name: GetArchive :one
SELECT * FROM archive_schema.archive WHERE id = $1 LIMIT 1;

-- name: GetArchiveByName :one
SELECT filename FROM archive_schema.archive WHERE filename=name;

-- name: InsertFile :exec
INSERT INTO archive_schema.archive (filename, editorial, file, thumbnail_image)
VALUES($1, $2, $3, $4);

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
