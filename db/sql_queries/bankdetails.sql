-- name: CreateBankDetails :one
INSERT INTO bank_details (bank_name, account_name, account_number, branch, swift_code, address, country, currency, account_type, company_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;



-- name: GetBankDetailsByAccountNumber :one
SELECT * FROM bank_details WHERE account_number = $1 LIMIT 1 ;


-- name: ListBankDetailsByBankName :many
SELECT * FROM bank_details WHERE bank_name = $1;

-- name: GetBankInfoByID :one
SELECT * FROM bank_details WHERE id = $1 LIMIT 1;

-- name: ListBanks :many
SELECT * FROM bank_details;