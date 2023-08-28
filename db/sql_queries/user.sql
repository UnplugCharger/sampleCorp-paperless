
-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email,
  department
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: DeleteUser :exec
DELETE FROM "users" WHERE id = $1;

-- name: GetUserByUserNameOrEmail :one
SELECT * FROM "users" WHERE username = $1 OR email = $2;

-- name: GetUsers :many
SELECT * FROM "users" ORDER BY id DESC LIMIT $1 OFFSET $2;

-- name: GetUserById :one
SELECT * FROM "users" WHERE id = $1;


-- name: UpdateUser :one
UPDATE "users" SET username = $1, email = $2 WHERE id = $3 RETURNING *;