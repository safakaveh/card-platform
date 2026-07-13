-- name: UpsertLaserData :one
INSERT INTO "card-platform".laser_data (
    uuid_card, 
    side, 
    row_no, 
    content_type, 
    content
) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (uuid_card, side, row_no) 
DO UPDATE SET 
    content_type = EXCLUDED.content_type,
    content = EXCLUDED.content,
    created_at = NOW()
RETURNING *;

-- name: GetLaserDataByCard :many
SELECT * FROM "card-platform".laser_data
WHERE uuid_card = $1;
