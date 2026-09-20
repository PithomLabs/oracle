package http

import (
	"net/http"

	coordinator "github.com/PithomLabs/oracle/coordinator"
)

// Server is the HTTP server for the Coordinator.
type Server struct {
	handler *Handler
}

// NewServer creates a new HTTP server.
func NewServer(c *coordinator.Coordinator) *Server {
	return &Server{
		handler: NewHandler(c),
	}
}

// Handler returns the http.Handler for the server.
func (s *Server) Handler() http.Handler {
	return s.handler
}
