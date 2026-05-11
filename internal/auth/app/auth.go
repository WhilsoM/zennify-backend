package app

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/zennify/backend/internal/auth/ports"
)

func (s *Service) Register(ctx context.Context, username, password string) (userID string, createdAt time.Time, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("app: hash password: %w", err)
	}

	id, created, err := s.users.CreateUser(ctx, username, string(hash))
	if err != nil {
		return "", time.Time{}, err
	}

	return id, created, nil
}

func (s *Service) Login(ctx context.Context, username, password string) (accessToken, refreshToken string, err error) {
	u, err := s.users.UserByUsername(ctx, username)
	if err != nil {
		if err == ports.ErrUserNotFound {
			return "", "", ports.ErrInvalidCredentials
		}
		return "", "", fmt.Errorf("app: load user: %w", err)
	}

	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
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
		return "", "", fmt.Errorf("app: save refresh session: %w", err)
	}

	return access, refresh, nil
}

func (s *Service) RefreshTokens(ctx context.Context, refreshToken string) (accessToken, refreshTokenOut string, err error) {
	userIDClaim, usernameClaim, jti, err := s.tokens.parseRefresh(refreshToken)
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
		return "", "", fmt.Errorf("app: delete refresh jti: %w", err)
	}

	access, err := s.tokens.mintAccess(userIDClaim, usernameClaim)
	if err != nil {
		return "", "", fmt.Errorf("app: mint access token: %w", err)
	}

	newRefresh, _, err := s.tokens.mintRefresh(userIDClaim, usernameClaim)
	if err != nil {
		return "", "", err
	}

	return access, newRefresh, nil
}
