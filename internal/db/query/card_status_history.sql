-- name: CreateCardStatusHistory :exec
INSERT INTO card_platform_card_status_history (
    uuid,
    uuid_card,
    status,
    created_at
) VALUES (?, ?, ?, ?);

-- name: GetCardStatusHistoryByUUID :one
SELECT
    uuid,
    uuid_card,
    status,
    created_at
FROM card_platform_card_status_history
WHERE uuid = ?
LIMIT 1;

-- name: ListCardStatusHistoryByCardUUID :many
SELECT
    uuid,
    uuid_card,
    status,
    created_at
FROM card_platform_card_status_history
WHERE uuid_card = ?
ORDER BY created_at ASC;

-- name: GetLatestCardStatusByCardUUID :one
SELECT
    uuid,
    uuid_card,
    status,
    created_at
FROM card_platform_card_status_history
WHERE uuid_card = ?
ORDER BY created_at DESC
LIMIT 1;

-- name: DeleteCardStatusHistoryByUUID :exec
DELETE FROM card_platform_card_status_history
WHERE uuid = ?;

-- name: DeleteCardStatusHistoriesByCardUUID :exec
DELETE FROM card_platform_card_status_history
WHERE uuid_card = ?;

-- name: CountCardStatusHistoryByCardUUID :one
SELECT COUNT(*) AS count
FROM card_platform_card_status_history
WHERE uuid_card = ?;
