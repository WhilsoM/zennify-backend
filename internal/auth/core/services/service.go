package services

import (
	"fmt"
	"time"

	"github.com/zennify/backend/internal/auth/core/ports"
)

type Service struct {
	users      ports.UserClient
	sessions   ports.RefreshSessionStore
	tokens     *tokenIssuer
	refreshTTL time.Duration
}

func NewService(users ports.UserClient, sessions ports.RefreshSessionStore, jwtSecret []byte, accessTTL, refreshTTL time.Duration) (*Service, error) {
	ti, err := newTokenIssuer(jwtSecret, accessTTL, refreshTTL)
	if err != nil {
		return nil, fmt.Errorf("auth services: token issuer: %w", err)
	}
	if refreshTTL <= 0 {
		refreshTTL = 7 * 24 * time.Hour
	}
	return &Service{
		users:      users,
		sessions:   sessions,
		tokens:     ti,
		refreshTTL: refreshTTL,
	}, nil
}
