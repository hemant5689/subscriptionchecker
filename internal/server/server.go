package server

import (
	"hej/internal/config"
	"hej/internal/router"
	"log"
	"net/http"
	"time"
)

// Server represents the API server
type Server struct {
	router *router.Router
	cfg    *config.Config
}

// NewServer creates a new server instance
func NewServer() *Server {
	cfg := config.LoadConfig()

	return &Server{
		router: router.NewRouter(cfg),
		cfg:    cfg,
	}
}

// Start initializes and starts the server
func (s *Server) Start() error {
	// Setup routes
	s.router.Setup()

	// Configure server
	server := &http.Server{
		Addr:         ":" + s.cfg.Port,
		Handler:      http.DefaultServeMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Server starting on port %s", s.cfg.Port)
	return server.ListenAndServe()
}
