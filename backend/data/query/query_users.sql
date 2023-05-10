-- name: FindUserByID :one
SELECT id, email, password, created_at, updated_at
FROM users
WHERE id = $1
LIMIT 1;

-- name: FindUserByEmail :one
SELECT id, email, password, created_at, updated_at
FROM users
WHERE email = $1
LIMIT 1;

