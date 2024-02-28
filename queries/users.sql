
-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: CheckUserExists :one
SELECT EXISTS (SELECT 1 FROM users WHERE email = $1) AS exists;

-- name: CreateUser :one
INSERT INTO users (username, password_hash, email, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;