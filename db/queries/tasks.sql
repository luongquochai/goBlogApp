-- name: CreateTask :one
INSERT INTO tasks (title, content, user_id)
VALUES ($1, $2, $3)
RETURNING id, title, content, user_id;

-- name: GetTaskByUserID :many
SELECT id, title, content, user_id
FROM tasks
WHERE user_id = $1;

-- name: UpdateTask :exec
UPDATE tasks
SET title = $2, content = $3
WHERE id = $1 AND user_id = $4;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = $1 AND user_id = $2;