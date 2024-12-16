package server

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/clholm/memory-gecko/youtube"
)

// handles requests to the /healthz endpoint
func handleHealthz(logger *log.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if err := encode(w, r, 200, "ok :)"); err != nil {
				http.NotFound(w, r)
			}
		},
	)
}

// handles requests to the index
func handleIndex(logger *log.Logger, config *Config) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// get project root directory
			projectRoot := getProjectRoot()

			// parse index.html
			tmpl, err := template.ParseFiles(filepath.Join(projectRoot, "web", "index.html"))
			// TODO: handle error lol
			if err != nil {
				logger.Printf("error parsing template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// execute template with video data
			err = tmpl.Execute(w, struct {
				Videos []youtube.SearchResult
			}{
				Videos: config.Videos,
			})
			if err != nil {
				logger.Printf("error executing template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		},
	)
}
