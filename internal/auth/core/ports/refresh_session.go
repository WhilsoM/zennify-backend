package ports

import (
	"context"
	"time"
)

type RefreshSessionStore interface {
	Ping(ctx context.Context) error
	SaveRefreshJTI(ctx context.Context, jti, userID string, expiresAt time.Time) error
	DeleteRefreshJTI(ctx context.Context, jti string) error
	UserIDForRefreshJTI(ctx context.Context, jti string) (userID string, err error)
}
