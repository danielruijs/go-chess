-- name: CreateUser :one
INSERT INTO users (username, display_name, hashed_password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: UserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);
