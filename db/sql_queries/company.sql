-- name: CreateCompany :one
INSERT INTO companies (name, initials, address)
VALUES ($1, $2, $3)
RETURNING *;



-- name: DeleteCompanyByName :exec
DELETE FROM companies
WHERE name = $1;

-- name: GetCompanyByID :one
SELECT * FROM companies
WHERE id = $1;

-- name: ListCompanies :many
SELECT * FROM companies LIMIT $1 OFFSET $2;