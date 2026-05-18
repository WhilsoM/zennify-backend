package services

import (
	"fmt"
	"time"

	"github.com/zennify/backend/internal/auth/core/ports"
)

type Service struct {
	users  ports.UserClient
	tokens *tokenIssuer
}

func NewService(users ports.UserClient, jwtSecret []byte, accessTTL, refreshTTL time.Duration) (*Service, error) {
	ti, err := newTokenIssuer(jwtSecret, accessTTL, refreshTTL)
	if err != nil {
		return nil, fmt.Errorf("auth services: token issuer: %w", err)
	}
	return &Service{
		users:  users,
		tokens: ti,
	}, nil
}
