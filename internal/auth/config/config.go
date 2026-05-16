package config

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	GRPCAddr            string        `env:"AUTH_GRPC_ADDR" envDefault:":50051"`
	UserServiceGRPCAddr string        `env:"AUTH_USER_SERVICE_GRPC_ADDR" envDefault:":50052"`
	JWTSecret           string        `env:"JWT_SECRET,required"`
	AccessTTL           time.Duration `env:"AUTH_ACCESS_TTL" envDefault:"15m"`
	RefreshTTL          time.Duration `env:"AUTH_REFRESH_TTL" envDefault:"168h"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("auth config: %w", err)
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("auth config: parse env: %w", err)
	}

	return &cfg, nil
}
