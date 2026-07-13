-- name: UpsertMagnetData :one
INSERT INTO "card-platform".magnet_data (
    uuid_card, 
    track_no, 
    content
) VALUES ($1, $2, $3)
ON CONFLICT (uuid_card, track_no) 
DO UPDATE SET 
    content = EXCLUDED.content,
    created_at = NOW()
RETURNING *;

-- name: GetMagnetDataByCard :many
SELECT * FROM "card-platform".magnet_data
WHERE uuid_card = $1;
