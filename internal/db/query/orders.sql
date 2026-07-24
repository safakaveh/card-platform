-- name: CreateOrder :exec
INSERT INTO orders (
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
FROM orders
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
FROM orders
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
FROM orders
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
FROM orders
WHERE status = ?
ORDER BY created_at DESC;

-- name: UpdateOrderStatus :exec
UPDATE orders
SET
    status = ?,
    updated_at = ?
WHERE uuid = ?;

-- name: UpdateOrderInfo :exec
UPDATE orders
SET
    order_name = ?,
    status = ?,
    description = ?,
    order_date = ?,
    updated_at = ?
WHERE uuid = ?;

-- name: UpdateOrderDescription :exec
UPDATE orders
SET
    description = ?,
    updated_at = ?
WHERE uuid = ?;

-- name: DeleteOrderByUUID :exec
DELETE FROM orders
WHERE uuid = ?;

-- name: CountOrders :one
SELECT COUNT(*) AS count
FROM orders;

-- name: CountOrdersByStatus :one
SELECT COUNT(*) AS count
FROM orders
WHERE status = ?;
