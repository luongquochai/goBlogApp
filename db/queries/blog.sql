-- Insert a new user
-- name: CreateUser :one
INSERT INTO users (username, password_hash)
VALUES ($1, $2)
RETURNING id, username, created_at;

-- Get a user by username
-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- Insert a new post
-- name: CreatePost :one
INSERT INTO posts (user_id, title, content)
VALUES ($1, $2, $3)
RETURNING *;

-- Get a post by ID
-- name: GetPost :one
SELECT * FROM posts
WHERE id = $1;

-- List all posts
-- name: ListPosts :many
SELECT * FROM posts
ORDER BY created_at DESC;

-- Update a post
-- name: UpdatePost :exec
UPDATE posts
SET title = $2, content = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- Delete a post
-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;