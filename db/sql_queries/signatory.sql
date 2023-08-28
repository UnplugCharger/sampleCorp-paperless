-- name: CreateSignatory :one
INSERT INTO signatories (name, title)
VALUES ($1, $2)
RETURNING *;



-- name: DeleteSignatoryByName :exec
DELETE FROM signatories
WHERE name = $1;

-- name: GetSignatoryById :one
SELECT * FROM signatories
WHERE id = $1;

-- name: ListSignatories :many
SELECT * FROM signatories LIMIT $1 OFFSET $2;