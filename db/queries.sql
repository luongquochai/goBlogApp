-- name: CreateUser :one
INSERT INTO users (username, email, hashed_password)
VALUES ($1, $2, $3)
RETURNING id, username, email, created_at, updated_at;

-- name: GetUserByUsername :one
SELECT id, username, email, hashed_password, created_at, updated_at
FROM users
WHERE username = $1;

-- name: UpdatePassword :exec
UPDATE users SET hashed_password = $2
WHERE username = $1;

-- name: CreatePost :one
INSERT INTO posts (title, content, author_id)
VALUES ($1, $2, $3)
RETURNING id, title, content, author_id, created_at, updated_at;

-- name: GetPostByID :one
SELECT id, title, content, author_id, created_at, updated_at
FROM posts
WHERE id = $1;

-- name: ListPosts :many
SELECT id, title, content, author_id, created_at updated_at
FROM posts
ORDER BY created_at DESC;

-- name: UpdatePostByAuthor :exec
UPDATE posts
SET title = $2, content = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND author_id = $4;

-- name: DeletePostByAuthor :exec
DELETE FROM posts
WHERE id = $1 AND author_id = $2;