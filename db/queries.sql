-- name: CreateUser :one
INSERT INTO users (username, hashed_password) VALUES ($1, $2)
RETURNING id, username;

-- name: GetUserByUsername :one
SELECT id, username, hashed_password FROM users
WHERE username = $1;

-- name: UpdatePassword :exec
UPDATE users SET hashed_password = $2
WHERE username = $1;