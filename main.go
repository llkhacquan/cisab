package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/llkhacquan/cisab/pkg/api"
	"github.com/llkhacquan/cisab/pkg/config"
	"github.com/llkhacquan/cisab/pkg/dbctx"
	"github.com/llkhacquan/cisab/pkg/repo"
	"github.com/llkhacquan/cisab/pkg/service"
	"github.com/llkhacquan/cisab/pkg/utils/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize logger
	appLogger := logger.NewDefault()

	// Load configuration
	appConfig, err := config.LoadConfigFromEnv()
	if err != nil {
		appLogger.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize database connection
	dbConfig := appConfig.Database
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		appLogger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	// Initialize repositories
	userRepo := repo.NewUserRepoImpl(dbctx.Get)
	taskRepo := repo.NewTaskRepoImpl(dbctx.Get)

	// Initialize services
	userService := service.NewUserService(userRepo, appConfig.JWT)
	taskService := service.NewTaskService(taskRepo, userRepo)

	// Create API server with services
	apiServer := api.NewServer(userService, taskService, userRepo, appLogger, db, appConfig.JWT.Secret)

	// Configure the HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", appConfig.Server.Port),
		Handler:      apiServer.Router(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the server
	appLogger.Info("server starting", "port", appConfig.Server.Port, "environment", appConfig.Environment)
	if err := server.ListenAndServe(); err != nil {
		appLogger.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
