package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

// following guide/tips at https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/

// config type contains server config info
type Config struct {
	Host   string
	Port   string
	APIKey string
}

// server constructor
func NewServer(
	logger *log.Logger,
	config *Config,
) (http.Handler, error) {
	mux := http.NewServeMux()
	err := addRoutes(
		mux,
		logger,
		config,
	)
	var handler http.Handler = mux
	if err != nil {
		return handler, err
	}
	// handler = logging.NewLoggingMiddleware(logger, handler)
	// handler = logging.NewGoogleTraceIDMiddleware(logger, handler)
	// handler = checkAuthHeaders(handler)
	return handler, nil
}

func serv(logger *log.Logger, config *Config, ctx context.Context, stderr io.Writer) error {
	// create server instance
	srv, err := NewServer(
		logger,
		config,
	)
	if err != nil {
		return err
	}
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: srv,
	}
	// goroutine calls ListenAndServe
	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(stderr, "error listening and serving: %s\n", err)
		}
	}()
	// waitgroup that waits on server shutdown goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	// server shutdown goroutine
	// shuts down when ctx tells it to
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}

// essentially the "main" function of this server
func Run(
	ctx context.Context,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	host, port, apiKey string,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// create Config object
	config := Config{
		Host:   host,
		Port:   port,
		APIKey: apiKey,
	}

	// create log object
	loggo := log.New(stdout, "memory-gecko:", log.LstdFlags)
	// call serv
	err := serv(loggo, &config, ctx, stderr)
	if err != nil {
		return err
	}
	// ...
	return nil
}
