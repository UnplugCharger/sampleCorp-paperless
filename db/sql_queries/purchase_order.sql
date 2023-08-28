-- name: CreatePurchaseOrder :one
INSERT INTO purchase_orders (attn, company_id, address, signatory_id, quotation_id, sent_or_received)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetPurchaseOrder :one
SELECT * FROM purchase_orders WHERE id = $1;

-- name: ListPurchaseOrders :many
SELECT * FROM purchase_orders  LIMIT $1 OFFSET $2;

-- name: UpdatePurchaseOrder :one
UPDATE purchase_orders
SET
    attn = COALESCE($1, attn),
    company_id = COALESCE($2, company_id),
    address = COALESCE($3, address),
    signatory_id = COALESCE($4, signatory_id),
    sent_or_received = COALESCE($5, sent_or_received)
WHERE
        id = $6
RETURNING *;

-- name: ApprovePurchaseOrder :one
UPDATE purchase_orders
SET
    po_status = $1,
    approved_by = $2
WHERE
        id = $3
RETURNING *;


