package main

import (
	"fmt"
	"log"

	"go.uber.org/zap"

	grpcstore "github.com/zennify/backend/internal/gateway/adapters/grpc"
	http "github.com/zennify/backend/internal/gateway/adapters/http"
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
		zap.String("auth", cfg.GrpcAddrs.AuthGRPCAddr),
		zap.String("user", cfg.GrpcAddrs.UserGRPCAddr),
	)

	authClient, userClient, err := initGRPCClients(cfg, logger)
	if err != nil {
		return fmt.Errorf("init grpc clients: %w", err)
	}

	svc := services.NewService(authClient, userClient, []byte(cfg.JWTSecret), cfg.RequestTimeout)
	router := http.NewRouter(svc, logger)

	return httpserver.Run(cfg.HTTPAddr, "api-gateway", cfg.ShutdownTimeout, router)
}

func initGRPCClients(cfg *config.Config, logger *zap.Logger) (*grpcstore.AuthClient, *grpcstore.UserClient, error) {
	authConn, err := grpcstore.NewAuthConn(cfg.GrpcAddrs.AuthGRPCAddr)
	if err != nil {
		return nil, nil, fmt.Errorf("auth conn: %w", err)
	}
	defer func() {
		if err := authConn.Close(); err != nil {
			logger.Error("close auth grpc connection", zap.Error(err))
		}
	}()

	userConn, err := grpcstore.NewUserConn(cfg.GrpcAddrs.UserGRPCAddr)
	if err != nil {
		return nil, nil, fmt.Errorf("user conn: %w", err)
	}
	defer func() {
		if err := userConn.Close(); err != nil {
			logger.Error("close user grpc connection", zap.Error(err))
		}
	}()

	authClient := grpcstore.NewAuthClient(authConn)
	userClient := grpcstore.NewUserClient(userConn)

	return authClient, userClient, nil
}
