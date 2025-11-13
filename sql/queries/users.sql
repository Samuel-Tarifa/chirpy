-- name: CreateUser :one
INSERT INTO users (email,hashed_password)
VALUES (
    $1,
    $2
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email=$1;

-- name: UpdateUser :one
UPDATE users
SET email=$2,hashed_password=$3,updated_at=$4
WHERE id=$1
RETURNING *;

-- name: SetChirpyRedUser :one
UPDATE users
SET is_chirpy_red=$2
WHERE id=$1
RETURNING *;