package memory

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/zennify/backend/internal/auth/ports"
)

var _ ports.UserClient = (*UserClientStub)(nil)

type userRow struct {
	id           string
	username     string
	passwordHash string
	createdAt    time.Time
}

type UserClientStub struct {
	mu         sync.RWMutex
	byUsername map[string]userRow
}

func NewUserClientStub() *UserClientStub {
	return &UserClientStub{
		byUsername: make(map[string]userRow),
	}
}

func (s *UserClientStub) Ping(ctx context.Context) error {
	_ = ctx
	return nil
}

func (s *UserClientStub) CreateUser(ctx context.Context, username, passwordHash string) (string, time.Time, error) {
	_ = ctx
	key := normalizeUsername(username)

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.byUsername[key]; exists {
		return "", time.Time{}, ports.ErrUsernameTaken
	}

	id := uuid.NewString()
	now := time.Now().UTC()
	s.byUsername[key] = userRow{
		id:           id,
		username:     username,
		passwordHash: passwordHash,
		createdAt:    now,
	}
	return id, now, nil
}

func (s *UserClientStub) UserByUsername(ctx context.Context, username string) (ports.User, error) {
	_ = ctx
	key := normalizeUsername(username)

	s.mu.RLock()
	defer s.mu.RUnlock()

	row, ok := s.byUsername[key]
	if !ok {
		return ports.User{}, ports.ErrUserNotFound
	}
	return ports.User{
		ID:           row.id,
		Username:     row.username,
		PasswordHash: row.passwordHash,
		CreatedAt:    row.createdAt,
	}, nil
}

func normalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}
