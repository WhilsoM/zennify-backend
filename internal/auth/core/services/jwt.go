package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	claimTypeAccess  = "access"
	claimTypeRefresh = "refresh"
)

type tokenIssuer struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func newTokenIssuer(jwtSecret []byte, accessTTL, refreshTTL time.Duration) (*tokenIssuer, error) {
	if len(jwtSecret) == 0 {
		return nil, fmt.Errorf("auth services: jwt secret is empty")
	}
	if accessTTL <= 0 {
		accessTTL = 15 * time.Minute
	}
	if refreshTTL <= 0 {
		refreshTTL = 7 * 24 * time.Hour
	}
	return &tokenIssuer{
		secret:     jwtSecret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}, nil
}

func (t *tokenIssuer) mintAccess(userID, username string) (string, error) {
	now := time.Now().UTC()
	exp := now.Add(t.accessTTL)

	claims := jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"typ":      claimTypeAccess,
		"iat":      now.Unix(),
		"exp":      exp.Unix(),
	}

	signed, signErr := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(t.secret)
	if signErr != nil {
		return "", fmt.Errorf("auth services: sign access token: %w", signErr)
	}
	return signed, nil
}

func (t *tokenIssuer) mintRefresh(userID, username string) (token string, jti string, err error) {
	jti = uuid.NewString()
	now := time.Now().UTC()
	exp := now.Add(t.refreshTTL)

	claims := jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"typ":      claimTypeRefresh,
		"jti":      jti,
		"iat":      now.Unix(),
		"exp":      exp.Unix(),
	}

	signed, signErr := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(t.secret)
	if signErr != nil {
		return "", "", fmt.Errorf("auth services: sign refresh token: %w", signErr)
	}
	return signed, jti, nil
}

func (t *tokenIssuer) parseRefresh(token string) (userID, username, jti string, err error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("auth services: unexpected signing method: %v", token.Header["alg"])
		}
		return t.secret, nil
	})
	if err != nil || !parsed.Valid {
		return "", "", "", fmt.Errorf("auth services: parse refresh token: %w", err)
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", "", fmt.Errorf("auth services: refresh claims type")
	}
	if typ, _ := claims["typ"].(string); typ != claimTypeRefresh {
		return "", "", "", fmt.Errorf("auth services: refresh token typ")
	}

	usernameClaim, _ := claims["username"].(string)
	sub, _ := claims["sub"].(string)
	jtiClaim, _ := claims["jti"].(string)
	if sub == "" || jtiClaim == "" {
		return "", "", "", fmt.Errorf("auth services: refresh token missing sub/jti")
	}
	return sub, usernameClaim, jtiClaim, nil
}
