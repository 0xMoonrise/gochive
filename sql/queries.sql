-- name: GetArchive :one
SELECT * FROM archive_schema.archive WHERE id = $1 LIMIT 1;

-- name: GetArchiveByName :one
SELECT filename FROM archive_schema.archive WHERE filename=name;

-- name: InsertFile :exec
INSERT INTO archive_schema.archive (filename, editorial, file)
VALUES($1, $2, $3);
