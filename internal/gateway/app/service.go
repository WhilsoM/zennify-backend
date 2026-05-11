package app

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sony/gobreaker"

	"github.com/zennify/backend/internal/gateway/ports"
)

type Service struct {
	auth           ports.AuthClient
	breaker        *gobreaker.CircuitBreaker
	jwtSecret      []byte
	requestTimeout time.Duration
}

type AccessClaims struct {
	UserID   string
	Username string
}

func NewService(auth ports.AuthClient, jwtSecret []byte, requestTimeout time.Duration) *Service {
	if requestTimeout <= 0 {
		requestTimeout = 5 * time.Second
	}
	return &Service{
		auth: auth,
		breaker: gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "gateway-auth-service",
			Timeout: 20 * time.Second,
		}),
		jwtSecret:      jwtSecret,
		requestTimeout: requestTimeout,
	}
}

func (s *Service) ValidateAccessToken(token string) (*AccessClaims, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil || !parsed.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	typ, _ := claims["typ"].(string)
	if typ != "access" {
		return nil, fmt.Errorf("invalid token type")
	}
	userID, _ := claims["sub"].(string)
	username, _ := claims["username"].(string)
	if userID == "" {
		return nil, fmt.Errorf("missing sub")
	}
	return &AccessClaims{UserID: userID, Username: username}, nil
}
