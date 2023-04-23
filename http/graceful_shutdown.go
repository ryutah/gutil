package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// StartServerwithGracefulShutdown start http server with graceful shutdown
func StartServerwithGracefulShutdown(addr string, h http.Handler, cleanFunc func()) error {
	srv := &http.Server{
		Handler: h,
		Addr:    addr,
	}

	errChan := make(chan error)
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			errChan <- fmt.Errorf("failed to listen server: %w ", err)
		}
	}()

	go func() {
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
		defer stop()
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if cleanFunc != nil {
			cleanFunc()
		}

		if err := srv.Shutdown(ctx); err != nil {
			errChan <- fmt.Errorf("failed to shutdown server: %w", err)
		}
		close(errChan)
	}()

	return <-errChan
}
