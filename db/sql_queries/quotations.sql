-- name: CreateQuotation :one
INSERT INTO quotations ( attn, company_id, site, validity, warranty, payment_terms, delivery_terms, signatory_id, status, sent_or_received)
VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: ListQuotations :many
SELECT * FROM quotations
LIMIT $1 OFFSET $2 ;


-- name: UpdateQuotation :one
UPDATE quotations
SET attn = COALESCE($1, attn),
    company_id = COALESCE($2, company_id),
    site = COALESCE($3, site),
    validity = COALESCE($4, validity),
    warranty = COALESCE($5, warranty),
    payment_terms = COALESCE($6, payment_terms),
    delivery_terms = COALESCE($7, delivery_terms),
    signatory_id = COALESCE($8, signatory_id),
    status = COALESCE($9, status),
    sent_or_received = COALESCE($10, sent_or_received)
WHERE id = $11
RETURNING *;

-- name: GetQuotationByID :one
SELECT * FROM quotations
WHERE id = $1 ;