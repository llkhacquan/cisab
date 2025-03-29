package api

import (
	"encoding/json"
	"net/http"

	"github.com/llkhacquan/knovel-assignment/pkg/models"
)

// Response is a generic response structure
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// HealthCheckHandler handles health check requests
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status:  "success",
		Message: "API is running",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonResponse); err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}

// getUsers returns mock user data
func getUsers() []models.User {
	return []models.User{
		{ID: 1, Username: "johndoe", Email: "john@example.com", Role: "admin"},
		{ID: 2, Username: "janedoe", Email: "jane@example.com", Role: "user"},
		{ID: 3, Username: "bobsmith", Email: "bob@example.com", Role: "user"},
	}
}

// GetUsersHandler handles GET requests for users
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Get mock data
	users := getUsers()

	response := Response{
		Status:  "success",
		Message: "Users retrieved successfully",
		Data:    users,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonResponse); err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}
