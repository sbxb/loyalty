package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/sbxb/loyalty/internal/logger"
)

type HTTPServer struct {
	srv             *http.Server
	shutdownTimeout time.Duration
}

var ErrServerStartFailed = errors.New("HTTPServer failed to start")

// NewHTTPServer creates a new server
func NewHTTPServer(address string, router http.Handler) (*HTTPServer, error) {
	// Set more reasonable timeouts than the default ones
	server := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  8 * time.Second,
		WriteTimeout: 8 * time.Second,
		IdleTimeout:  36 * time.Second,
	}

	return &HTTPServer{
		srv:             server,
		shutdownTimeout: 3 * time.Second,
	}, nil
}

func (s *HTTPServer) Start(ctx context.Context) error {
	if s.srv == nil {
		return ErrServerStartFailed
	}

	logger.Info("HTTPServer ready to start")
	if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
		s.srv = nil
		logger.Errorf("HTTPServer ListenAndServe() failed: %v", err)
		return ErrServerStartFailed
	}

	logger.Info("HTTPServer is being gracefully stopped")
	return nil
}

func (s *HTTPServer) Close() {
	if s.srv == nil {
		return
	}
	logger.Info("Trying to gracefully stop HTTPServer")
	// Perform server shutdown with a default maximum timeout of 3 seconds
	timeoutCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.srv.Shutdown(timeoutCtx); err != nil {
		// Error from closing listeners, or context timeout:
		logger.Errorf("HTTPServer Shutdown() failed: %v", err)
	} else {
		logger.Info("HTTPServer has been gracefully stopped")
	}

	s.srv = nil
}
