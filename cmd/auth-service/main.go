package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	authgrpc "github.com/zennify/backend/internal/auth/adapters/grpc"
	authconfig "github.com/zennify/backend/internal/auth/config"
	"github.com/zennify/backend/internal/auth/core/services"
	"github.com/zennify/backend/internal/shared/grpcserver"
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

	cfg, err := authconfig.LoadConfig()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users, err := authgrpc.NewUserClient(ctx, cfg.UserServiceGRPCAddr)
	if err != nil {
		return fmt.Errorf("user service client: %w", err)
	}
	defer func() {
		if closeErr := users.Close(); closeErr != nil {
			logger.Error("user service client close", zap.Error(closeErr))
		}
	}()

	svc, err := services.NewService(users, []byte(cfg.JWTSecret), cfg.AccessTTL, cfg.RefreshTTL)
	if err != nil {
		return fmt.Errorf("services: %w", err)
	}

	return grpcserver.Run(cfg.GRPCAddr, "auth-service", 10*time.Second, func(s *grpc.Server) {
		authgrpc.Register(s, svc)
		reflection.Register(s)
	})
}
