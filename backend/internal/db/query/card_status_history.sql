-- name: AddCardStatusHistory :one
INSERT INTO "card-platform".card_status_history (
    uuid_card, 
    status
) VALUES ($1, $2)
RETURNING *;

-- name: GetCardHistory :many
SELECT * FROM "card-platform".card_status_history
WHERE uuid_card = $1
ORDER BY created_at DESC;
