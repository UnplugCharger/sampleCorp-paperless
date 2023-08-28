-- name: CreateRole :one
-- description: Create a role
INSERT INTO roles (name, description) values ($1, $2) returning *;

-- name: DeleteRole :exec
-- description: Delete a role
DELETE FROM roles WHERE id = $1 returning *;


-- name: GetRole :one
-- description: Get a role
SELECT * FROM roles WHERE id = $1;

-- name: UpdateRole :one
-- description: Update a role
UPDATE roles SET name = $1, description = $2 WHERE id = $3 returning *;

-- name: ListRoles :many
-- description: List roles
SELECT * FROM roles ORDER BY id ASC LIMIT $1 OFFSET $2;