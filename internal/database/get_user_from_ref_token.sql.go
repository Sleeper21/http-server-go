// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: get_user_from_ref_token.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const getUserFromRefreshToken = `-- name: GetUserFromRefreshToken :one
SELECT users.email, is_chirpy_red, refresh_tokens.token, refresh_tokens.created_at, refresh_tokens.updated_at, refresh_tokens.user_id, refresh_tokens.expires_at, refresh_tokens.revoked_at
FROM users
JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1
`

type GetUserFromRefreshTokenRow struct {
	Email       string
	IsChirpyRed bool
	Token       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      uuid.UUID
	ExpiresAt   time.Time
	RevokedAt   sql.NullTime
}

func (q *Queries) GetUserFromRefreshToken(ctx context.Context, token string) (GetUserFromRefreshTokenRow, error) {
	row := q.db.QueryRowContext(ctx, getUserFromRefreshToken, token)
	var i GetUserFromRefreshTokenRow
	err := row.Scan(
		&i.Email,
		&i.IsChirpyRed,
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}
