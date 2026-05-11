package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func Run(addr, serviceName string, shutdownTimeout time.Duration, handler http.Handler) error {
	if shutdownTimeout <= 0 {
		shutdownTimeout = 10 * time.Second
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serveErrCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErrCh <- err
		}
		close(serveErrCh)
	}()

	log.Printf("%s HTTP listening on %s", serviceName, addr)

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("http shutdown: %w", err)
		}
		return nil
	case serveErr := <-serveErrCh:
		if serveErr != nil {
			return fmt.Errorf("http serve: %w", serveErr)
		}
		return nil
	}
}
