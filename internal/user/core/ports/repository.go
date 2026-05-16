package ports

import (
	"context"

	"github.com/zennify/backend/internal/user/core/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByUsername(ctx context.Context, username string) (domain.User, error)
}
