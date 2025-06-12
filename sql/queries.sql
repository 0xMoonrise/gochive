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
	editorial
FROM archive_schema.archive 
ORDER BY id
LIMIT $1
OFFSET $2; 

-- name: GetCountArchive :one
SELECT
	count(id)
FROM archive_schema.archive;

-- name: GetThumbnails :many
SELECT filename, thumbnail_image FROM archive_schema.archive;
