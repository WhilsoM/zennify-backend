package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserRecord struct {
	ID           uuid.UUID
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

type UserRepository interface {
	Create(ctx context.Context, user UserRecord) error
	GetByUsername(ctx context.Context, username string) (UserRecord, error)
}
