-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, hashed_password, email)
VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- get user by email address

-- name: GetUserByEmail :one

SELECT users.* from users
WHERE users.email = $1;