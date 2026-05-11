package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/zennify/backend/internal/auth/ports"
)

var _ ports.RefreshSessionStore = (*RefreshSessionStore)(nil)

type refreshRow struct {
	userID    string
	expiresAt time.Time
}

type RefreshSessionStore struct {
	mu    sync.Mutex
	byJTI map[string]refreshRow
}

func NewRefreshSessionStore() *RefreshSessionStore {
	return &RefreshSessionStore{
		byJTI: make(map[string]refreshRow),
	}
}

func (s *RefreshSessionStore) Ping(ctx context.Context) error {
	_ = ctx
	return nil
}

func (s *RefreshSessionStore) SaveRefreshJTI(ctx context.Context, jti, userID string, expiresAt time.Time) error {
	_ = ctx
	if jti == "" || userID == "" {
		return fmt.Errorf("memory: save refresh jti: %w", ports.ErrInvalidRefreshToken)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.byJTI[jti] = refreshRow{userID: userID, expiresAt: expiresAt.UTC()}
	return nil
}

func (s *RefreshSessionStore) DeleteRefreshJTI(ctx context.Context, jti string) error {
	_ = ctx

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.byJTI, jti)
	return nil
}

func (s *RefreshSessionStore) UserIDForRefreshJTI(ctx context.Context, jti string) (string, error) {
	_ = ctx

	s.mu.Lock()
	defer s.mu.Unlock()

	row, ok := s.byJTI[jti]
	if !ok {
		return "", ports.ErrInvalidRefreshToken
	}
	if time.Now().UTC().After(row.expiresAt) {
		delete(s.byJTI, jti)
		return "", ports.ErrInvalidRefreshToken
	}
	return row.userID, nil
}
