-- name: CreateQuotationItem :one
INSERT INTO quotation_items (
    description,
    uom,
    qty,
    lead_time,
    item_price,
    disc,
    unit_price,
    net_price,
    currency,
    quotation_id
) VALUES (
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
         ) RETURNING *;

-- name: GetQuotationItemsByQuotationID :many
SELECT *
FROM quotation_items
WHERE quotation_id = $1;


-- name: UpdateQuotationItem :one
UPDATE quotation_items
SET
    description = COALESCE($1, description),
    uom = COALESCE($2, uom),
    qty = COALESCE($3, qty),
    lead_time = COALESCE($4, lead_time),
    item_price = COALESCE($5, item_price),
    disc = COALESCE($6, disc),
    unit_price = COALESCE($7, unit_price),
    net_price = COALESCE($8, net_price),
    currency = COALESCE($9, currency),
    quotation_id = COALESCE($10, quotation_id)
WHERE id = $11
RETURNING *;


-- name: ListQuotationItemsByQuotationID :many
SELECT * FROM quotation_items WHERE quotation_id = $1 ;
