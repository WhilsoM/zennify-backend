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

	access, err := s.tokens.mintAccess(u.ID, u.Username)
	if err != nil {
		return "", "", err
	}

	refresh, jti, err := s.tokens.mintRefresh(u.ID, u.Username)
	if err != nil {
		return "", "", err
	}

	expiresAt := time.Now().UTC().Add(s.refreshTTL)
	if err := s.sessions.SaveRefreshJTI(ctx, jti, u.ID, expiresAt); err != nil {
		return "", "", fmt.Errorf("auth services: save refresh session: %w", err)
	}

	return access, refresh, nil
}

func (s *Service) RefreshTokens(ctx context.Context, req ports.RefreshTokensRequest) (accessToken, refreshTokenOut string, err error) {
	userIDClaim, usernameClaim, jti, err := s.tokens.parseRefresh(req.RefreshToken)
	if err != nil {
		return "", "", ports.ErrInvalidRefreshToken
	}

	storedUserID, err := s.sessions.UserIDForRefreshJTI(ctx, jti)
	if err != nil {
		return "", "", ports.ErrInvalidRefreshToken
	}
	if storedUserID != userIDClaim {
		return "", "", ports.ErrInvalidRefreshToken
	}

	if err := s.sessions.DeleteRefreshJTI(ctx, jti); err != nil {
		return "", "", fmt.Errorf("auth services: delete refresh jti: %w", err)
	}

	access, err := s.tokens.mintAccess(userIDClaim, usernameClaim)
	if err != nil {
		return "", "", fmt.Errorf("auth services: mint access token: %w", err)
	}

	newRefresh, _, err := s.tokens.mintRefresh(userIDClaim, usernameClaim)
	if err != nil {
		return "", "", err
	}

	return access, newRefresh, nil
}
