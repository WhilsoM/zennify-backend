package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zennify/backend/internal/user/ports"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

const createUserSQL = `
INSERT INTO "user" (id, username, password_hash, created_at)
VALUES ($1, $2, $3, $4)
`

func (r *UserRepository) Create(ctx context.Context, user ports.UserRecord) error {
	_, err := r.pool.Exec(ctx, createUserSQL, user.ID, user.Username, user.PasswordHash, user.CreatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return ports.ErrUsernameTaken
		}
		return fmt.Errorf("postgres user repository: create: %w", err)
	}
	return nil
}

const getUserByUsernameSQL = `
SELECT id, username, password_hash, created_at
FROM "user"
WHERE username = $1
`

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (ports.UserRecord, error) {
	var user ports.UserRecord
	err := r.pool.QueryRow(ctx, getUserByUsernameSQL, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ports.UserRecord{}, ports.ErrUserNotFound
		}
		return ports.UserRecord{}, fmt.Errorf("postgres user repository: get by username: %w", err)
	}
	return user, nil
}

func isUniqueViolation(err error) bool {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		return pgErr.Code == "23505"
	}
	return false
}
