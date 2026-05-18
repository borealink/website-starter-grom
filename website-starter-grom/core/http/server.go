package http

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"website-starter/core/logging"
)

// Server wraps the standard HTTP server
// and provides graceful shutdown support.
type Server struct {
	httpServer *http.Server
	logger     *logging.Logger
}

// NewServer creates and configures a new HTTP server instance.
func NewServer(
	addr string,
	handler http.Handler,
	logger *logging.Logger,
) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		logger: logger,
	}
}

// Start launches the HTTP server in a separate goroutine.
func (server *Server) Start() {
	go func() {

		// Log server startup message.
		server.logger.Get().Println(
			"Server running on",
			server.httpServer.Addr,
		)

		// Start listening for incoming HTTP requests.
		err := server.httpServer.ListenAndServe()

		// Ignore the normal server closed error during shutdown.
		if err != nil && !errors.Is(
			err,
			http.ErrServerClosed,
		) {
			server.logger.Get().Fatal(
				err,
			)
		}
	}()
}

// Shutdown gracefully stops the server.
func (server *Server) Shutdown(
	ctx context.Context,
) error {

	// Log shutdown start.
	server.logger.Get().Println(
		"Shutting down server...",
	)

	// Attempt graceful shutdown.
	err := server.httpServer.Shutdown(
		ctx,
	)

	if err != nil {
		return err
	}

	// Log successful shutdown.
	server.logger.Get().Println(
		"Server stopped",
	)

	return nil
}

// WaitForShutdown waits for an interrupt signal
// and gracefully shuts down the server.
func (server *Server) WaitForShutdown(
	timeout time.Duration,
) {

	// Create signal channel.
	stop := make(
		chan os.Signal,
		1,
	)

	// Listen for interrupt and termination signals.
	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	// Block until a signal is received.
	<-stop

	// Create shutdown timeout context.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		timeout,
	)

	defer cancel()

	// Shutdown the server gracefully.
	err := server.Shutdown(
		ctx,
	)

	if err != nil {
		server.logger.Get().Fatal(
			err,
		)
	}
}
