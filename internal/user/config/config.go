package config

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	GRPCAddr string `env:"USER_GRPC_ADDR" envDefault:":50052"`
	Database DatabaseConfig
}

type DatabaseConfig struct {
	URL               string        `env:"USER_DATABASE_URL,required"`
	MaxConns          int32         `env:"USER_DB_MAX_CONNS" envDefault:"25"`
	MinConns          int32         `env:"USER_DB_MIN_CONNS" envDefault:"2"`
	MaxConnLifetime   time.Duration `env:"USER_DB_MAX_CONN_LIFETIME" envDefault:"1h"`
	MaxConnIdleTime   time.Duration `env:"USER_DB_MAX_CONN_IDLE_TIME" envDefault:"30m"`
	HealthCheckPeriod time.Duration `env:"USER_DB_HEALTH_CHECK_PERIOD" envDefault:"1m"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("user config: %w", err)
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("user config: parse env: %w", err)
	}

	return &cfg, nil
}
