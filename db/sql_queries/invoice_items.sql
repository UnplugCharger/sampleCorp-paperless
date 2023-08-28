 -- name: CreateInvoiceItem :one
 INSERT INTO invoice_items (description, uom, qty, unit_price, net_price, currency, invoice_id)
 VALUES ($1, $2, $3, $4, $5, $6, $7)
 RETURNING *;


-- name: GetInvoiceItemsByInvoiceID :many
 SELECT id, description, uom, qty, unit_price, net_price, currency, invoice_id
 FROM invoice_items
 WHERE invoice_id = $1;


-- name: ListInvoiceItemsByInvoiceID :many
 SELECT id, description, uom, qty, unit_price, net_price, currency, invoice_id
 FROM invoice_items
 WHERE invoice_id = $1;