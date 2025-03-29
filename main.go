package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/llkhacquan/knovel-assignment/pkg/api"
	"github.com/llkhacquan/knovel-assignment/pkg/config"
	"github.com/llkhacquan/knovel-assignment/pkg/utils/logger"
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

	// Create router with all routes and middleware
	router := api.NewRouter()

	// Configure the HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", appConfig.Server.Port),
		Handler:      router,
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
