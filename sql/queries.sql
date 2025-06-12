-- name: GetArchive :one
SELECT id, filename FROM archive_schema.archive WHERE id = $1;

-- name: GetArchiveByName :one
SELECT filename FROM archive_schema.archive WHERE filename=name;
