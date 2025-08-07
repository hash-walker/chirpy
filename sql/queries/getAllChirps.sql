-- name: GetAllChirps :many

SELECT chirps.* from chirps
ORDER BY chirps.created_at ASC;

-- name: GetChirp :one
SELECT chirps.* from chirps
WHERE chirps.id = $1;