-- name: CreateOrder :exec
INSERT INTO card_platform_orders (
    uuid,
    order_name,
    status,
    description,
    order_date,
    created_at,
    updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetOrderByUUID :one
SELECT
    uuid,
    order_name,
    status,
    description,
    order_date,
    created_at,
    updated_at
FROM card_platform_orders
WHERE uuid = ?
LIMIT 1;

-- name: GetOrderByName :one
SELECT
    uuid,
    order_name,
    status,
    description,
    order_date,
    created_at,
    updated_at
FROM card_platform_orders
WHERE order_name = ?
LIMIT 1;

-- name: ListOrders :many
SELECT
    uuid,
    order_name,
    status,
    description,
    order_date,
    created_at,
    updated_at
FROM card_platform_orders
ORDER BY created_at DESC;

-- name: ListOrdersByStatus :many
SELECT
    uuid,
    order_name,
    status,
    description,
    order_date,
    created_at,
    updated_at
FROM card_platform_orders
WHERE status = ?
ORDER BY created_at DESC;

-- name: UpdateOrderStatus :exec
UPDATE card_platform_orders
SET
    status = ?,
    updated_at = ?
WHERE uuid = ?;

-- name: UpdateOrderInfo :exec
UPDATE card_platform_orders
SET
    order_name = ?,
    status = ?,
    description = ?,
    order_date = ?,
    updated_at = ?
WHERE uuid = ?;

-- name: UpdateOrderDescription :exec
UPDATE card_platform_orders
SET
    description = ?,
    updated_at = ?
WHERE uuid = ?;

-- name: DeleteOrderByUUID :exec
DELETE FROM card_platform_orders
WHERE uuid = ?;

-- name: CountOrders :one
SELECT COUNT(*) AS count
FROM card_platform_orders;

-- name: CountOrdersByStatus :one
SELECT COUNT(*) AS count
FROM card_platform_orders
WHERE status = ?;
