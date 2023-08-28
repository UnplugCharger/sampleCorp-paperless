-- name: CreateInvoice :one
INSERT INTO invoices (purchase_order_number, attn, company_id, site, amount_due, bank_details, signatory_id, sent_or_received)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;



-- name: GetInvoiceById :one
SELECT * FROM invoices WHERE id = $1;

-- name: GetInvoicesByPurchaseOrderNumber :many
SELECT * FROM invoices WHERE purchase_order_number = $1;



-- name: ListInvoices :many
SELECT * FROM invoices
LIMIT $1 OFFSET $2;

-- name: ListUserCompanyInvoices :many
SELECT * FROM invoices WHERE company_id = $1
LIMIT $2 OFFSET $3;

-- name: ListSiteInvoices :many
SELECT * FROM invoices WHERE site = $1
LIMIT $2 OFFSET $3;