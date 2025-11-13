-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token=$1;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token,user_id,expires_at)
VALUES(
  $1,
  $2,
  $3
) 
RETURNING *;

-- name: UpdateRefreshToken :one
UPDATE refresh_tokens
SET expires_at=$2,updated_at=$3
WHERE token=$1
RETURNING *;