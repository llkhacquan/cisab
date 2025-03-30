package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/llkhacquan/knovel-assignment/pkg/repo"
	"github.com/llkhacquan/knovel-assignment/pkg/service"
	"github.com/llkhacquan/knovel-assignment/pkg/utils/logger"
	"gorm.io/gorm"
)

// Server represents the HTTP server
type Server struct {
	router      *mux.Router
	logger      *logger.Logger
	userService service.UserService
	taskService service.TaskService
	userRepo    repo.UserRepo
	gormDB      *gorm.DB
	jwtSecret   string
}

// NewServer creates a new HTTP server
func NewServer(userService service.UserService, taskService service.TaskService, userRepo repo.UserRepo, log *logger.Logger, gormDB *gorm.DB, jwtSecret string) *Server {
	server := &Server{
		router:      mux.NewRouter(),
		logger:      log,
		userService: userService,
		taskService: taskService,
		userRepo:    userRepo,
		gormDB:      gormDB,
		jwtSecret:   jwtSecret,
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

	// Add a global OPTIONS handler for CORS preflight requests
	s.router.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The CORS middleware will set the necessary headers and return 200 OK
		w.WriteHeader(http.StatusOK)
	})

	// Non-API endpoints (health check)
	healthEndpoint := []Endpoint{
		{
			Method:  http.MethodGet,
			Path:    "/health",
			Handler: s.HealthCheckHandler,
		},
	}

	// Register health endpoint directly on main router
	RegisterEndpoints(s.router, healthEndpoint, s.logger)

	// Create API router with middlewares
	apiRouter := s.router.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(DBTransactionMiddleware(s.logger, s.gormDB))
	apiRouter.Use(AuthMiddleware(s.logger, s.userRepo, s.jwtSecret))

	// All API endpoints
	apiEndpoints := []Endpoint{
		// Public endpoints (authentication is handled by the middleware)
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: s.LoginHandler,
		},
		// User endpoints
		{
			Method:  http.MethodGet,
			Path:    "/users/me",
			Handler: s.GetMeHandler,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/all",
			Handler: s.GetUsersHandler,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{id}",
			Handler: s.GetUserByIDHandler,
		},
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: s.CreateUserHandler,
		},
		// Task endpoints
		{
			Method:  http.MethodPost,
			Path:    "/tasks",
			Handler: s.CreateTaskHandler,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/tasks/{id}/status",
			Handler: s.UpdateTaskStatusHandler,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/tasks/{id}/assign",
			Handler: s.AssignTaskHandler,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks/assigned",
			Handler: s.GetAssignedTasksHandler,
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks",
			Handler: s.GetTasksHandler,
		},
		// Employee summary endpoint
		{
			Method:  http.MethodGet,
			Path:    "/employee-summary",
			Handler: s.GetEmployeeTaskSummaryHandler,
		},
	}

	// Register all API endpoints
	RegisterEndpoints(apiRouter, apiEndpoints, s.logger)
}

// Router returns the server's router
func (s *Server) Router() *mux.Router {
	return s.router
}
