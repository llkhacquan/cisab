package api

import (
	"github.com/gorilla/mux"
	"github.com/llkhacquan/knovel-assignment/pkg/utils/logger"
)

// NewRouter creates and configures a new router with all routes and middleware
func NewRouter() *mux.Router {
	router := mux.NewRouter()

	// Create logger
	log := logger.NewDefault()

	// Apply global middleware
	router.Use(LoggingMiddleware(log))
	router.Use(CorsMiddleware)

	// Public routes
	router.HandleFunc("/health", HealthCheckHandler).Methods("GET")

	// API routes with auth middleware
	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(AuthMiddleware(log))

	// Register protected routes
	apiRouter.HandleFunc("/users", GetUsersHandler).Methods("GET")

	return router
}
