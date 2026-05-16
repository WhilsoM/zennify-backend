package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/zennify/backend/internal/shared/grpcserver"
	sharedpostgres "github.com/zennify/backend/internal/shared/postgres"
	userpostgres "github.com/zennify/backend/internal/user/adapters/db/postgres"
	usergrpc "github.com/zennify/backend/internal/user/adapters/grpc"
	userConfig "github.com/zennify/backend/internal/user/config"
	"github.com/zennify/backend/internal/user/core/services"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("run: %v", err)
	}
}

func run() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("logger: %w", err)
	}
	defer func() {
		if syncErr := logger.Sync(); syncErr != nil {
			log.Printf("logger sync: %v", syncErr)
		}
	}()

	cfg, err := userConfig.LoadConfig()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := sharedpostgres.NewPool(ctx, sharedpostgres.PoolConfig{
		URL:               cfg.Database.URL,
		MaxConns:          cfg.Database.MaxConns,
		MinConns:          cfg.Database.MinConns,
		MaxConnLifetime:   cfg.Database.MaxConnLifetime,
		MaxConnIdleTime:   cfg.Database.MaxConnIdleTime,
		HealthCheckPeriod: cfg.Database.HealthCheckPeriod,
	})
	if err != nil {
		return fmt.Errorf("database: %w", err)
	}
	defer func() {
		pool.Close()
		logger.Info("database pool closed")
	}()

	userRepo := userpostgres.NewUserRepository(pool)
	svc := services.NewService(userRepo)

	return grpcserver.Run(cfg.GRPCAddr, "user-service", 10*time.Second, func(s *grpc.Server) {
		usergrpc.Register(s, svc)
		reflection.Register(s)
	})
}
