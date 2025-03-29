package api

import (
	"encoding/json"
	"net/http"
)

// HealthCheckHandler handles health check requests
func (s *Server) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status:  "success",
		Message: "API is running",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		s.logger.Error("failed to marshal health check response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonResponse); err != nil {
		s.logger.Error("failed to write health check response", "error", err)
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}

// GetUsersHandler handles GET requests for users
func (s *Server) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Get users from the service layer
	users, err := s.userService.GetUsers(r.Context())
	if err != nil {
		s.logger.Error("failed to get users", "error", err)
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	response := Response{
		Status:  "success",
		Message: "Users retrieved successfully",
		Data:    users,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		s.logger.Error("failed to marshal users response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonResponse); err != nil {
		s.logger.Error("failed to write users response", "error", err)
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}
