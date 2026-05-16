package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/zennify/backend/internal/auth/app"
	authconfig "github.com/zennify/backend/internal/auth/config"
	usergrpc "github.com/zennify/backend/internal/auth/store/grpc"
	"github.com/zennify/backend/internal/auth/store/memory"
	grpcapi "github.com/zennify/backend/internal/auth/transport/grpc"
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

	users, err := usergrpc.NewUserClient(ctx, cfg.UserServiceGRPCAddr)
	if err != nil {
		return fmt.Errorf("user service client: %w", err)
	}
	defer func() {
		if closeErr := users.Close(); closeErr != nil {
			logger.Error("user service client close", zap.Error(closeErr))
		}
	}()

	sessions := memory.NewRefreshSessionStore()
	svc, err := app.NewService(users, sessions, []byte(cfg.JWTSecret), cfg.AccessTTL, cfg.RefreshTTL)
	if err != nil {
		return fmt.Errorf("app: %w", err)
	}

	return grpcserver.Run(cfg.GRPCAddr, "auth-service", 10*time.Second, func(s *grpc.Server) {
		grpcapi.Register(s, svc)
		reflection.Register(s)
	})
}
