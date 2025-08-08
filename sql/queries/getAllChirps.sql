-- name: GetAllChirps :many

SELECT chirps.* from chirps
ORDER BY chirps.created_at;

-- name: GetChirp :one
SELECT chirps.* from chirps
WHERE chirps.id = $1;

-- name: GetChirpByAuthor :many
SELECT chirps.* from chirps
WHERE chirps.user_id = $1
ORDER BY chirps.created_at ASC;