-- name: CreateOrder :one
INSERT INTO "card-platform".orders (
    order_name, 
    description
) VALUES ($1, $2)
RETURNING *;

-- name: GetOrder :one
SELECT * FROM "card-platform".orders
WHERE uuid = $1 LIMIT 1;

-- name: ListOrders :many
SELECT * FROM "card-platform".orders
ORDER BY created_at DESC;

-- name: UpdateOrderStatus :one
UPDATE "card-platform".orders
SET status = $2
WHERE uuid = $1
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM "card-platform".orders
WHERE uuid = $1;
