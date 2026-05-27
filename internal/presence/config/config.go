package config

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	GRPCAddr string `env:"PRESENCE_GRPC_ADDR" envDefault:":50053"`
	NATSURL  string `env:"NATS_URL" envDefault:"nats://localhost:4222"`
	Redis    RedisConfig
}

type RedisConfig struct {
	RedisAddr                string        `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	RedisPassword            string        `env:"REDIS_PASSWORD" envDefault:""`
	RedisDB                  int           `env:"REDIS_DB" envDefault:"0"`
	RedisTTL                 time.Duration `env:"REDIS_TTL" envDefault:"15m"`
	RedisCleanupInterval     time.Duration `env:"REDIS_CLEANUP_INTERVAL" envDefault:"1m"`
	RedisCleanupBatchSize    int           `env:"REDIS_CLEANUP_BATCH_SIZE" envDefault:"1000"`
	RedisCleanupBatchTimeout time.Duration `env:"REDIS_CLEANUP_BATCH_TIMEOUT" envDefault:"10s"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("presence config: %w", err)
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("presence config: parse env: %w", err)
	}

	return &cfg, nil
}
