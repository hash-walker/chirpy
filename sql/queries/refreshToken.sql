-- name: CreateRefreshToken :one
INSERT INTO revoke_token (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: CheckToken :one
SELECT revoke_token.* from revoke_token
WHERE revoke_token.token = $1;

-- name: RevokeToken :exec
UPDATE revoke_token SET revoked_at = NOW(), updated_at = NOW() WHERE token = $1;

