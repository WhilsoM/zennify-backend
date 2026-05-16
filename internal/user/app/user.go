package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/zennify/backend/internal/user/ports"
)

func (s *Service) CreateUser(ctx context.Context, req ports.CreateUserRequest) (ports.CreateUserResponse, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return ports.CreateUserResponse{}, fmt.Errorf("user app: hash password: %w", err)
	}

	id := uuid.New()
	createdAt := time.Now().UTC()

	err = s.users.Create(ctx, ports.UserRecord{
		ID:           id,
		Username:     req.Username,
		PasswordHash: string(passwordHash),
		CreatedAt:    createdAt,
	})
	if err != nil {
		return ports.CreateUserResponse{}, err
	}

	return ports.CreateUserResponse{
		UserID:    id.String(),
		Username:  req.Username,
		CreatedAt: createdAt,
	}, nil
}

func (s *Service) GetUserByUsername(ctx context.Context, username string) (ports.UserRecord, error) {
	user, err := s.users.GetByUsername(ctx, username)
	if err != nil {
		return ports.UserRecord{}, err
	}
	return user, nil
}
