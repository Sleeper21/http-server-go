-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES ($1, now(), now(), $2, $3, NULL); -- the revoked_at should be NULL when the token is created

-- name: UpdateRefreshToken :exec
UPDATE refresh_tokens
SET token = $1, updated_at = now(), expires_at = $2, revoked_at = now()
WHERE token = $3;

-- name: GetRefreshTokenByToken :one
SELECT * FROM refresh_tokens WHERE token = $1;