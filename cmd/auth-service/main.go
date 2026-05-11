package main

import (
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/zennify/backend/internal/auth/app"
	authconfig "github.com/zennify/backend/internal/auth/config"
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
	cfg, err := authconfig.LoadConfig()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	users := memory.NewUserClientStub()
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
