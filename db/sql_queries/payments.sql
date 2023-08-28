-- name: CreatePaymentRequest :one
INSERT INTO payment_requests (
employee_id,
currency,
amount,
description,
status,
amount_in_words


) VALUES (
$1, $2, $3, $4, $5, $6
         ) RETURNING *;




-- name: UpdatePaymentRequest :one
UPDATE payment_requests
SET
    employee_id = COALESCE($1, employee_id),
    amount = COALESCE($2, amount),
    description = COALESCE($3, description),
    invoice_id = COALESCE($4, invoice_id),
    admin_id = COALESCE($5, admin_id)
WHERE request_id = $6 AND lower(status) ='pending'
RETURNING *;





-- name: DeletePaymentRequest :exec
DELETE FROM payment_requests
WHERE request_id = $1;



-- name: ApprovePaymentRequest :one
UPDATE payment_requests
SET
status = $1,
approval_date = $2,
admin_id = $3

WHERE request_id = $4 AND lower(status) ='pending'
RETURNING *;


-- name: ListEmployeePaymentRequests :many
SELECT *
FROM payment_requests
WHERE employee_id = $1
ORDER BY request_date DESC
LIMIT $2 OFFSET $3;

-- name: ListPaymentRequests :many
SELECT * FROM payment_requests ORDER BY request_date DESC LIMIT $1 OFFSET $2;

-- name: GetPaymentRequest :one
SELECT * FROM payment_requests WHERE request_id = $1;


-- name: ApprovePettyCash :one
UPDATE petty_cash
SET
status = $1,
approved_at = $2,
authorised_by = $3

WHERE transaction_id = $4 AND lower(status) ='pending'
RETURNING *;

-- name: CreatePettyCash :one
INSERT INTO petty_cash (employee_id, amount, description, transaction_date, folio, debit_account)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;



-- name: UpdatePettyCash :one
UPDATE petty_cash
SET
    amount = COALESCE($1, amount),
    description = COALESCE($2, description)
WHERE transaction_id = $3 AND lower(status) ='pending'
RETURNING *;


-- name: DeletePettyCash :exec
DELETE FROM petty_cash
WHERE transaction_id = $1;

-- name: ListEmployeePettyCash :many
SELECT * FROM petty_cash
WHERE employee_id = $1
ORDER BY transaction_date DESC
LIMIT $2 OFFSET $3 ;

-- name: ListPettyCash :many
SELECT * FROM petty_cash ORDER BY transaction_date LIMIT $1 OFFSET $2;

-- name: GetPettyCash :one
SELECT * FROM petty_cash WHERE transaction_id = $1;
