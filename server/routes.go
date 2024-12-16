package server

import (
	"log"
	"net/http"
	"path/filepath"
)

func addRoutes(
	mux *http.ServeMux,
	logger *log.Logger,
	config *Config,
) error {
	// get project root directory
	// _, b, _, _ := runtime.Caller(0)
	projectRoot := getProjectRoot()

	// serve static files from web directory
	webDir := http.Dir(filepath.Join(projectRoot, "web"))
	fs := http.FileServer(webDir)

	// add routes
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.Handle("/healthz", handleHealthz(logger))
	mux.Handle("/", handleIndex(logger, config))

	return nil
}
