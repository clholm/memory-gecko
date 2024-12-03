package server

import (
	"log"
	"net/http"
)

func addRoutes(
	mux *http.ServeMux,
	logger *log.Logger,
	config *Config,
) error {
	// add error handling logic for errors returned by handlers (lol)
	// mux.Handle("/api/v1/", handleRootGet(logger))
	mux.Handle("/healthz", handleHealthz(logger))
	mux.Handle("/", http.NotFoundHandler())

	return nil
}
