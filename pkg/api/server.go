package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/llkhacquan/knovel-assignment/pkg/service"
	"github.com/llkhacquan/knovel-assignment/pkg/utils/logger"
	"gorm.io/gorm"
)

// Server represents the HTTP server
type Server struct {
	router      *mux.Router
	logger      *logger.Logger
	userService service.UserService
	gormDB      *gorm.DB
}

// NewServer creates a new HTTP server
func NewServer(userService service.UserService, log *logger.Logger, gormDB *gorm.DB) *Server {
	server := &Server{
		router:      mux.NewRouter(),
		logger:      log,
		userService: userService,
		gormDB:      gormDB,
	}

	// Set up routes
	server.setupRoutes()

	return server
}

// setupRoutes configures all the routes for the server
func (s *Server) setupRoutes() {
	// Apply global middleware
	s.router.Use(LoggingMiddleware(s.logger))
	s.router.Use(CorsMiddleware)

	// Public endpoints
	publicEndpoints := []Endpoint{
		{
			Method:  http.MethodGet,
			Path:    "/health",
			Handler: s.HealthCheckHandler,
		},
	}

	// Register public endpoints
	RegisterEndpoints(s.router, publicEndpoints, s.logger)

	// API router with auth middleware
	apiRouter := s.router.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(DBTransactionMiddleware(s.logger, s.gormDB))
	apiRouter.Use(AuthMiddleware(s.logger))

	// User endpoints
	userEndpoints := []Endpoint{
		{
			Method:  http.MethodGet,
			Path:    "/users/{id}",
			Handler: s.GetUserByIDHandler,
		},
	}

	// Register user endpoints
	RegisterEndpoints(apiRouter, userEndpoints, s.logger)
}

// Router returns the server's router
func (s *Server) Router() *mux.Router {
	return s.router
}
