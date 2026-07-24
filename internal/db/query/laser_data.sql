-- name: CreateLaserData :exec
INSERT INTO laser_data (
    uuid,
    uuid_card,
    side,
    row_no,
    content_type,
    content,
    created_at
) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: UpsertLaserData :exec
INSERT INTO laser_data (
    uuid,
    uuid_card,
    side,
    row_no,
    content_type,
    content,
    created_at
) VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(uuid_card, side, row_no) DO UPDATE SET
    content_type = excluded.content_type,
    content = excluded.content;

-- name: GetLaserDataByUUID :one
SELECT
    uuid,
    uuid_card,
    side,
    row_no,
    content_type,
    content,
    created_at
FROM laser_data
WHERE uuid = ?
LIMIT 1;

-- name: GetLaserDataByCardSideRow :one
SELECT
    uuid,
    uuid_card,
    side,
    row_no,
    content_type,
    content,
    created_at
FROM laser_data
WHERE uuid_card = ?
  AND side = ?
  AND row_no = ?
LIMIT 1;

-- name: ListLaserDataByCardUUID :many
SELECT
    uuid,
    uuid_card,
    side,
    row_no,
    content_type,
    content,
    created_at
FROM laser_data
WHERE uuid_card = ?
ORDER BY side ASC, row_no ASC;

-- name: ListLaserDataByCardUUIDAndSide :many
SELECT
    uuid,
    uuid_card,
    side,
    row_no,
    content_type,
    content,
    created_at
FROM laser_data
WHERE uuid_card = ?
  AND side = ?
ORDER BY row_no ASC;

-- name: UpdateLaserDataContent :exec
UPDATE laser_data
SET
    content_type = ?,
    content = ?
WHERE uuid = ?;

-- name: DeleteLaserDataByUUID :exec
DELETE FROM laser_data
WHERE uuid = ?;

-- name: DeleteLaserDataByCardUUID :exec
DELETE FROM laser_data
WHERE uuid_card = ?;

-- name: DeleteLaserDataByCardSideRow :exec
DELETE FROM laser_data
WHERE uuid_card = ?
  AND side = ?
  AND row_no = ?;

-- name: CountLaserDataByCardUUID :one
SELECT COUNT(*) AS count
FROM laser_data
WHERE uuid_card = ?;
