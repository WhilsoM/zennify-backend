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

// Run grpc server with graceful shutdown
// @param addr - address to listen
// @param serviceName - name of the service
// @param shutdownTimeout - timeout for graceful shutdown
// @param configure - function to configure the grpc server
// @return error if any
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
