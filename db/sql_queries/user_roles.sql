-- name: CreateUserRoles :one
INSERT INTO user_roles (user_id, role_id)values ($1, $2) RETURNING *;


-- name: UpdateUserRoles :one
UPDATE user_roles SET role_id = $2 WHERE user_id = $1 RETURNING *;


-- name: DeleteUserRoles :exec
DELETE FROM user_roles WHERE user_id = $1;

-- name: DeleteUserRolesByRole :exec
DELETE FROM user_roles WHERE role_id = $1;

-- name: GetUserRoles :many
SELECT * FROM user_roles WHERE user_id = $1;


-- name: GetUserRole :one
SELECT * FROM user_roles WHERE user_id = $1 AND role_id = $2;

