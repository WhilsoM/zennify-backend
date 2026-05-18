// Package grpcserver runs a gRPC server with graceful shutdown on SIGINT/SIGTERM.
package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

// Run listens on addr, registers services via configure, and blocks until shutdown.
//
// configure is called with a new grpc.Server to register service implementations.
// On SIGINT or SIGTERM, GracefulStop is used; if it does not finish within
// shutdownTimeout, Stop forces termination. shutdownTimeout defaults to 10s when <= 0.
//
// serviceName is only used in the startup log line.
func Run(addr, serviceName string, shutdownTimeout time.Duration, configure func(*grpc.Server)) error {
	if shutdownTimeout <= 0 {
		shutdownTimeout = 10 * time.Second
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", addr, err)
	}
	defer func() {
		if err := lis.Close(); err != nil {
			log.Printf("grpc listener close: %v", err)
		}
	}()

	s := grpc.NewServer()
	configure(s)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serveErrCh := make(chan error, 1)
	go func() {
		if serveErr := s.Serve(lis); serveErr != nil && !errors.Is(serveErr, grpc.ErrServerStopped) {
			serveErrCh <- serveErr
		}
		close(serveErrCh)
	}()

	log.Printf("%s gRPC listening on %s", serviceName, addr)

	select {
	case <-ctx.Done():
		shutdownDone := make(chan struct{})
		go func() {
			s.GracefulStop()
			close(shutdownDone)
		}()

		select {
		case <-shutdownDone:
		case <-time.After(shutdownTimeout):
			s.Stop()
		}
		return nil
	case serveErr := <-serveErrCh:
		if serveErr != nil {
			return fmt.Errorf("serve: %w", serveErr)
		}
		return nil
	}
}
