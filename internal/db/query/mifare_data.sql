-- name: CreateMifareData :exec
INSERT INTO mifare_data (
    uuid,
    uuid_card,
    block_no,
    key_a,
    key_b,
    content,
    created_at
) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: UpsertMifareData :exec
INSERT INTO mifare_data (
    uuid,
    uuid_card,
    block_no,
    key_a,
    key_b,
    content,
    created_at
) VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(uuid_card, block_no) DO UPDATE SET
    key_a = excluded.key_a,
    key_b = excluded.key_b,
    content = excluded.content;

-- name: GetMifareDataByUUID :one
SELECT
    uuid,
    uuid_card,
    block_no,
    key_a,
    key_b,
    content,
    created_at
FROM mifare_data
WHERE uuid = ?
LIMIT 1;

-- name: GetMifareDataByCardBlock :one
SELECT
    uuid,
    uuid_card,
    block_no,
    key_a,
    key_b,
    content,
    created_at
FROM mifare_data
WHERE uuid_card = ?
  AND block_no = ?
LIMIT 1;

-- name: ListMifareDataByCardUUID :many
SELECT
    uuid,
    uuid_card,
    block_no,
    key_a,
    key_b,
    content,
    created_at
FROM mifare_data
WHERE uuid_card = ?
ORDER BY block_no ASC;

-- name: UpdateMifareDataContent :exec
UPDATE mifare_data
SET
    key_a = ?,
    key_b = ?,
    content = ?
WHERE uuid = ?;

-- name: DeleteMifareDataByUUID :exec
DELETE FROM mifare_data
WHERE uuid = ?;

-- name: DeleteMifareDataByCardUUID :exec
DELETE FROM mifare_data
WHERE uuid_card = ?;

-- name: DeleteMifareDataByCardBlock :exec
DELETE FROM mifare_data
WHERE uuid_card = ?
  AND block_no = ?;

-- name: CountMifareDataByCardUUID :one
SELECT COUNT(*) AS count
FROM mifare_data
WHERE uuid_card = ?;
