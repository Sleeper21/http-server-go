-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (gen_random_uuid(), now(), now(), $1, $2)
RETURNING *;

-- The :one at the end of the query name tells SQLC that we expect to get back a single row (the created user).

-- name: GetUserByID :one
SELECT id, created_at, updated_at, email FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET email = $1, hashed_password = $2, updated_at = now()
WHERE id = $3
RETURNING id, created_at, updated_at, email;

