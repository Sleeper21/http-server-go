-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (gen_random_uuid(), now(), now(), $1)
RETURNING *;

-- The :one at the end of the query name tells SQLC that we expect to get back a single row (the created user).

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

