package api

import (
	"github.com/gorilla/mux"
	"github.com/llkhacquan/knovel-assignment/pkg/service"
	"github.com/llkhacquan/knovel-assignment/pkg/utils/logger"
)

// Server represents the HTTP server
type Server struct {
	router      *mux.Router
	logger      *logger.Logger
	userService service.UserService
}

// NewServer creates a new HTTP server
func NewServer(userService service.UserService, log *logger.Logger) *Server {
	server := &Server{
		router:      mux.NewRouter(),
		logger:      log,
		userService: userService,
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

	// Public routes
	s.router.HandleFunc("/health", s.HealthCheckHandler).Methods("GET")

	// API routes with auth middleware
	apiRouter := s.router.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(AuthMiddleware(s.logger))

	// Register protected routes
	apiRouter.HandleFunc("/users", s.GetUsersHandler).Methods("GET")
}

// Router returns the server's router
func (s *Server) Router() *mux.Router {
	return s.router
}
