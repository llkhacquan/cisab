package api

import (
	"net/http"
)

// HealthCheckHandler handles health check requests
func (s *Server) HealthCheckHandler(r *http.Request) (interface{}, error) {
	return map[string]string{"status": "ok"}, nil
}
