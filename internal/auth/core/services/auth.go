package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/zennify/backend/internal/auth/core/ports"
)

func (s *Service) Register(ctx context.Context, req ports.RegisterRequest) (userID string, createdAt time.Time, err error) {
	id, created, err := s.users.CreateUser(ctx, req.Username, req.Password)
	if err != nil {
		return "", time.Time{}, err
	}

	return id, created, nil
}

func (s *Service) Login(ctx context.Context, req ports.LoginRequest) (accessToken, refreshToken string, err error) {
	u, err := s.users.UserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, ports.ErrUserNotFound) {
			return "", "", ports.ErrInvalidCredentials
		}
		return "", "", fmt.Errorf("auth services: load user: %w", err)
	}

	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
		return "", "", ports.ErrInvalidCredentials
	}

	return s.issueTokenPair(u.ID, u.Username)
}

func (s *Service) RefreshTokens(ctx context.Context, req ports.RefreshTokensRequest) (accessToken, refreshTokenOut string, err error) {
	_ = ctx
	userID, username, err := s.tokens.parseRefresh(req.RefreshToken)
	if err != nil {
		return "", "", ports.ErrInvalidRefreshToken
	}

	return s.issueTokenPair(userID, username)
}

func (s *Service) issueTokenPair(userID, username string) (accessToken, refreshToken string, err error) {
	access, err := s.tokens.mintAccess(userID, username)
	if err != nil {
		return "", "", err
	}

	refresh, err := s.tokens.mintRefresh(userID, username)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
