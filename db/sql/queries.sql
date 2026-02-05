-- name: GetArchive :one
SELECT * FROM archive WHERE id = ? LIMIT 1;

-- name: GetArchiveByName :one
SELECT * FROM archive WHERE filename = ? LIMIT 1;

-- name: GetArchiveById :one
SELECT filename FROM archive WHERE id = ?;

-- name: InsertFile :one
INSERT INTO archive(
  filename,
  editorial
)
VALUES(?, ?)
RETURNING id;

-- name: GetArchivePage :many
SELECT
  id,
  filename,
  editorial,
  favorite
FROM archive
ORDER BY favorite DESC, id
LIMIT ?
OFFSET ?;

-- name: GetCountArchive :one
SELECT
  count(id)
FROM archive;

-- name: SearchArchive :many
SELECT
  id,
  filename,
  editorial,
  favorite
FROM archive
WHERE filename LIKE '%' || ? || '%'
ORDER BY favorite DESC, id
LIMIT ?
OFFSET ?;

-- name: GetCountSearch :one
SELECT
  count(id)
FROM archive
WHERE filename LIKE '%' || ? || '%';

-- name: SetFavorite :exec
UPDATE archive
SET favorite = ?
WHERE id = ?;

-- name: SetEditFile :exec
UPDATE
  archive
SET
  filename = ?,
  editorial = ?
WHERE id = ?;
