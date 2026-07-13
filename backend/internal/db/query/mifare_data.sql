-- name: UpsertMifareData :one
INSERT INTO "card-platform".mifare_data (
    uuid_card, 
    block_no, 
    key_a, 
    key_b, 
    content
) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (uuid_card, block_no) 
DO UPDATE SET 
    key_a = EXCLUDED.key_a,
    key_b = EXCLUDED.key_b,
    content = EXCLUDED.content,
    created_at = NOW()
RETURNING *;

-- name: GetMifareDataByCard :many
SELECT * FROM "card-platform".mifare_data
WHERE uuid_card = $1;
