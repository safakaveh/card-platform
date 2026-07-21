-- name: CreateMagnetData :exec
INSERT INTO card_platform_magnet_data (
    uuid,
    uuid_card,
    track_no,
    content,
    created_at
) VALUES (?, ?, ?, ?, ?);

-- name: UpsertMagnetData :exec
INSERT INTO card_platform_magnet_data (
    uuid,
    uuid_card,
    track_no,
    content,
    created_at
) VALUES (?, ?, ?, ?, ?)
ON CONFLICT(uuid_card, track_no) DO UPDATE SET
    content = excluded.content;

-- name: GetMagnetDataByUUID :one
SELECT
    uuid,
    uuid_card,
    track_no,
    content,
    created_at
FROM card_platform_magnet_data
WHERE uuid = ?
LIMIT 1;

-- name: GetMagnetDataByCardTrack :one
SELECT
    uuid,
    uuid_card,
    track_no,
    content,
    created_at
FROM card_platform_magnet_data
WHERE uuid_card = ?
  AND track_no = ?
LIMIT 1;

-- name: ListMagnetDataByCardUUID :many
SELECT
    uuid,
    uuid_card,
    track_no,
    content,
    created_at
FROM card_platform_magnet_data
WHERE uuid_card = ?
ORDER BY track_no ASC;

-- name: UpdateMagnetDataContent :exec
UPDATE card_platform_magnet_data
SET
    content = ?
WHERE uuid = ?;

-- name: DeleteMagnetDataByUUID :exec
DELETE FROM card_platform_magnet_data
WHERE uuid = ?;

-- name: DeleteMagnetDataByCardUUID :exec
DELETE FROM card_platform_magnet_data
WHERE uuid_card = ?;

-- name: DeleteMagnetDataByCardTrack :exec
DELETE FROM card_platform_magnet_data
WHERE uuid_card = ?
  AND track_no = ?;

-- name: CountMagnetDataByCardUUID :one
SELECT COUNT(*) AS count
FROM card_platform_magnet_data
WHERE uuid_card = ?;
