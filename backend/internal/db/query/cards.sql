-- name: CreateCard :one
INSERT INTO "card-platform".cards (
    uuid_order, 
    has_laser, 
    has_magnet, 
    has_mifare_uid, 
    has_mifare, 
    has_java_card, 
    has_press, 
    has_temperature, 
    description
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetCard :one
SELECT * FROM "card-platform".cards
WHERE uuid = $1 LIMIT 1;

-- name: ListCardsByOrder :many
SELECT * FROM "card-platform".cards
WHERE uuid_order = $1
ORDER BY created_at ASC;

-- name: UpdateCardDoneStatus :one
UPDATE "card-platform".cards
SET is_done = $2
WHERE uuid = $1
RETURNING *;

-- name: DeleteCard :exec
DELETE FROM "card-platform".cards
WHERE uuid = $1;
