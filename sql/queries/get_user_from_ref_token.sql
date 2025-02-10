-- name: GetUserFromRefreshToken :one
SELECT users.email, refresh_tokens.*
FROM users
JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1;