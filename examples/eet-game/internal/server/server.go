// Package server wires together the HTTP routes and WebSocket endpoint.
package server

import (
	"fmt"
	"log"
	"net/http"
)

// Config holds server configuration options.
type Config struct {
	Port    int
	StaticDir string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Port:      8080,
		StaticDir: "web/static",
	}
}

// Run starts the HTTP server and blocks until it exits.
func Run(cfg Config) error {
	mux := http.NewServeMux()

	// WebSocket endpoint
	mux.HandleFunc("/ws", wsHandler)

	// Static files (HTML / JS / CSS)
	fs := http.FileServer(http.Dir(cfg.StaticDir))
	mux.Handle("/", fs)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("🚀 EET Space Invaders running on http://localhost%s", addr)
	return http.ListenAndServe(addr, mux)
}
