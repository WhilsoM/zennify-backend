package main

import (
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/zennify/backend/internal/shared/grpcserver"
	"github.com/zennify/backend/internal/user/app"
	userConfig "github.com/zennify/backend/internal/user/config"
	grpcapi "github.com/zennify/backend/internal/user/transport/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("run: %v", err)
	}
}

func run() error {
	cfg, err := userConfig.LoadConfig()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	svc := app.NewService()

	return grpcserver.Run(cfg.GRPCAddr, "user-service", 10*time.Second, func(s *grpc.Server) {
		grpcapi.Register(s, svc)
		reflection.Register(s)
	})
}
