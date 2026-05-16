package ports

import (
	"context"
	"time"

	"github.com/zennify/backend/internal/auth/core/domain"
)

type UserClient interface {
	Ping(ctx context.Context) error
	CreateUser(ctx context.Context, username, password string) (userID string, createdAt time.Time, err error)
	UserByUsername(ctx context.Context, username string) (domain.User, error)
}
