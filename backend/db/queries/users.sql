-- name: CreateUser :one
INSERT INTO users (username, display_name, hashed_password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;
