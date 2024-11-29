package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

// following guide/tips at https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/

// server constructor
func NewServer(
	logger *Logger,
	config *Config,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		Logger,
		Config,
	)
	var handler http.Handler = mux
	// handler = logging.NewLoggingMiddleware(logger, handler)
	// handler = logging.NewGoogleTraceIDMiddleware(logger, handler)
	// handler = checkAuthHeaders(handler)
	return handler
}

func addRoutes error (
	mux *http.ServeMux,
	logger *logging.Logger,
	config Config,
) {
	// add error handling logic for errors returned by handlers (lol)
	mux.Handle("/api/v1/", handleRootGet(logger))
	mux.HandleFunc("/healthz", handleHealthz(logger))
	mux.Handle("/", http.NotFoundHandler())

	return nil
}

func serv() error {
	// create server instance
	srv := NewServer(
		logger,
		config,
	)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: srv,
	}
	// goroutine calls ListenAndServe
	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
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
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}

// essentially the "main" function of this server
func run(
	ctx context.Context,
	stdin io.Reader,
	stdout io.Writer,
	args []string,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// ...
}

// TODO: put this in the function that starts the server 
// func main() {
// 	ctx := context.Background()
// 	if err := run(ctx, os.Stdout, os.Args); err != nil {
// 		fmt.Fprintf(os.Stderr, "%s\n", err)
// 		os.Exit(1)
// 	}
// }