package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// encode and decode helper functions
func encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

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
