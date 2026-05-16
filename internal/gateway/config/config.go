package config

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr string `env:"GATEWAY_HTTP_ADDR" envDefault:":8080"`
	ServicesAddresses
	JWTSecret       string        `env:"JWT_SECRET,required"`
	RequestTimeout  time.Duration `env:"GATEWAY_REQUEST_TIMEOUT" envDefault:"5s"`
	ShutdownTimeout time.Duration `env:"GATEWAY_SHUTDOWN_TIMEOUT" envDefault:"10s"`
}

type ServicesAddresses struct {
	AuthGRPCAddr string `env:"AUTH_GRPC_ADDR" envDefault:":50051"`
	UserGRPCAddr string `env:"USER_GRPC_ADDR"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("gateway config: %w", err)
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("gateway config: parse env: %w", err)
	}

	return &cfg, nil
}
