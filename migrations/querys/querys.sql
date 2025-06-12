-- name: GetArchive :one
SELECT id, filename FROM archive_schema.archive WHERE id = $1;
