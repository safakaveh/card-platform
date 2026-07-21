-- name: CreateCard :exec
INSERT INTO card_platform_cards (
    uuid,
    uuid_order,
    has_laser,
    has_magnet,
    has_mifare_uid,
    has_mifare,
    has_java_card,
    has_press,
    has_temperature,
    is_done,
    description,
    created_at,
    updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetCardByUUID :one
SELECT
    uuid,
    uuid_order,
    has_laser,
    has_magnet,
    has_mifare_uid,
    has_mifare,
    has_java_card,
    has_press,
    has_temperature,
    is_done,
    description,
    created_at,
    updated_at
FROM card_platform_cards
WHERE uuid = ?
LIMIT 1;

-- name: ListCardsByOrderUUID :many
SELECT
    uuid,
    uuid_order,
    has_laser,
    has_magnet,
    has_mifare_uid,
    has_mifare,
    has_java_card,
    has_press,
    has_temperature,
    is_done,
    description,
    created_at,
    updated_at
FROM card_platform_cards
WHERE uuid_order = ?
ORDER BY created_at DESC;

-- name: ListDoneCardsByOrderUUID :many
SELECT
    uuid,
    uuid_order,
    has_laser,
    has_magnet,
    has_mifare_uid,
    has_mifare,
    has_java_card,
    has_press,
    has_temperature,
    is_done,
    description,
    created_at,
    updated_at
FROM card_platform_cards
WHERE uuid_order = ?
  AND is_done = 1
ORDER BY created_at DESC;

-- name: ListPendingCardsByOrderUUID :many
SELECT
    uuid,
    uuid_order,
    has_laser,
    has_magnet,
    has_mifare_uid,
    has_mifare,
    has_java_card,
    has_press,
    has_temperature,
    is_done,
    description,
    created_at,
    updated_at
FROM card_platform_cards
WHERE uuid_order = ?
  AND is_done = 0
ORDER BY created_at DESC;

-- name: UpdateCardFlagsAndDescription :exec
UPDATE card_platform_cards
SET
    has_laser = ?,
    has_magnet = ?,
    has_mifare_uid = ?,
    has_mifare = ?,
    has_java_card = ?,
    has_press = ?,
    has_temperature = ?,
    is_done = ?,
    description = ?,
    updated_at = ?
WHERE uuid = ?;

-- name: UpdateCardDoneStatus :exec
UPDATE card_platform_cards
SET
    is_done = ?,
    updated_at = ?
WHERE uuid = ?;

-- name: UpdateCardDescription :exec
UPDATE card_platform_cards
SET
    description = ?,
    updated_at = ?
WHERE uuid = ?;

-- name: DeleteCardByUUID :exec
DELETE FROM card_platform_cards
WHERE uuid = ?;

-- name: CountCardsByOrderUUID :one
SELECT COUNT(*) AS count
FROM card_platform_cards
WHERE uuid_order = ?;

-- name: CountDoneCardsByOrderUUID :one
SELECT COUNT(*) AS count
FROM card_platform_cards
WHERE uuid_order = ?
  AND is_done = 1;
