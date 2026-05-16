package main

import (
	"fmt"
	"log"

	"go.uber.org/zap"

	grpcstore "github.com/zennify/backend/internal/gateway/adapters/grpc"
	httpapi "github.com/zennify/backend/internal/gateway/adapters/http"
	"github.com/zennify/backend/internal/gateway/config"
	"github.com/zennify/backend/internal/gateway/core/services"
	"github.com/zennify/backend/internal/shared/httpserver"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("run: %v", err)
	}
}

func run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("logger: %w", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("zap logger sync: %v", err)
		}
	}()

	logger.Info("gateway grpc targets",
		zap.String("auth", cfg.AuthGRPCAddr),
		zap.String("user", cfg.UserGRPCAddr),
	)

	authConn, err := grpcstore.NewAuthConn(cfg.AuthGRPCAddr)
	if err != nil {
		return fmt.Errorf("auth conn: %w", err)
	}
	defer func() {
		if err := authConn.Close(); err != nil {
			logger.Error("close auth grpc connection", zap.Error(err))
		}
	}()

	authClient := grpcstore.NewAuthClient(authConn)
	svc := services.NewService(authClient, []byte(cfg.JWTSecret), cfg.RequestTimeout)
	router := httpapi.NewRouter(svc, logger)

	return httpserver.Run(cfg.HTTPAddr, "api-gateway", cfg.ShutdownTimeout, router)
}
