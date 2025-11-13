-- name: CreateChirp :one
INSERT INTO chirps (body,user_id)
VALUES (
    $1,
    $2
)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps
WHERE user_id=$1 OR $1 IS NULL
ORDER BY created_at ASC;

-- name: GetChirp :one
SELECT * FROM chirps
WHERE id=$1;

-- name: DeleteChirp :one
DELETE FROM chirps
WHERE id=$1
RETURNING *;