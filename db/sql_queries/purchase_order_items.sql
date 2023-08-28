-- name: CreatePurchaseOrderItem :one
INSERT INTO purchase_order_items (
    description,
    part_no,
    uom,
    qty,
    item_price,
    discount,
    net_price,
    net_value,
    currency,
    purchase_order_id
) VALUES (
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
         ) RETURNING *;


-- name: GetPurchaseOrderItemsByPurchaseOrderID :many
SELECT *
FROM purchase_order_items
WHERE purchase_order_id = $1;


-- name: ListPurchaseOrderItemsByPurchaseOrderID :many
SELECT * FROM purchase_order_items WHERE purchase_order_id = $1 ;

-- name: UpdatePurchaseOrderItem :one
UPDATE purchase_order_items
SET
    description = COALESCE($1, description),
    part_no = COALESCE($2, part_no),
    uom = COALESCE($3, uom),
    qty = COALESCE($4, qty),
    item_price = COALESCE($5, item_price),
    discount = COALESCE($6, discount),
    net_price = COALESCE($7, net_price),
    net_value = COALESCE($8, net_value),
    currency = COALESCE($9, currency),
    purchase_order_id = COALESCE($10, purchase_order_id)
WHERE
        id = $11
RETURNING *;
