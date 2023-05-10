-- name: FindTaskByID :one
SELECT id, user_id, name, is_completed, created_at, updated_at
FROM tasks
WHERE id = $1
LIMIT 1;

-- name: FindTasksByUserID :many
SELECT id, user_id, name, is_completed, created_at, updated_at
FROM tasks
WHERE user_id = $1
ORDER BY updated_at DESC;

-- name: CreateTask :one
INSERT INTO tasks(id, user_id, name, is_completed, created_at, updated_at)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: UpdateTask :exec
UPDATE tasks
SET name = $2, is_completed = $3, updated_at = $4
WHERE id = $1;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1;
