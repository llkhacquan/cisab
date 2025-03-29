package testutil

import (
	"database/sql"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/llkhacquan/knovel-assignment/db"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// generateRandomDBName creates a random database name for testing
func generateRandomDBName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	const letters = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return "test_" + string(b)
}

// DBOption defines options for test database creation
type DBOption func(*gorm.Config)

// WithLogLevel sets the log level for the test database
func WithLogLevel(level logger.LogLevel) DBOption {
	return func(cfg *gorm.Config) {
		cfg.Logger = logger.Default.LogMode(level)
	}
}

// WithLogger sets a custom logger for the test database
func WithLogger(l logger.Interface) DBOption {
	return func(cfg *gorm.Config) {
		cfg.Logger = l
	}
}

// CreateTestDB creates a test database with migrations and returns a GORM DB instance
func CreateTestDB(t *testing.T, options ...DBOption) *gorm.DB {
	// Generate a random database name
	randomDBName := generateRandomDBName()
	t.Logf("Creating test database: %s", randomDBName)

	// Create the database and run migrations
	sqlDB, err := db.SetupTestDB(randomDBName)
	require.NoError(t, err, "failed to create test database")

	// Setup cleanup for when the test is done
	t.Cleanup(func() {
		cleanupTestDatabase(t, sqlDB, randomDBName)
	})

	// Create the GORM configuration with provided options
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Default to info level
	}

	// Apply custom options
	for _, option := range options {
		option(gormConfig)
	}

	// Create a GORM DB instance from the SQL DB
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), gormConfig)

	require.NoError(t, err, "failed to create GORM DB from SQL DB")
	return gormDB
}

// cleanupTestDatabase properly cleans up a test database
func cleanupTestDatabase(t *testing.T, db *sql.DB, dbName string) {
	// Get connection string to connect to postgres database
	dbConfig := struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}{
		Host:     "localhost",
		Port:     "5433",
		User:     "postgres",
		Password: "password",
		Name:     "postgres", // Connect to the default postgres database
	}

	// First close the existing connection
	if err := db.Close(); err != nil {
		t.Logf("Warning: failed to close test database connection: %v", err)
	}

	// Connect to postgres to drop the database
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name,
	)

	pgDB, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Logf("Warning: failed to connect to postgres database for cleanup: %v", err)
		return
	}
	defer pgDB.Close()

	// Terminate any active connections to the database
	// #nosec G201 - This is a test helper with controlled internal input, not user input
	terminateSQL := fmt.Sprintf(
		"SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '%s'",
		dbName,
	)

	if _, err := pgDB.Exec(terminateSQL); err != nil {
		t.Logf("Warning: failed to terminate connections to test database: %v", err)
	}

	// Drop the database
	// #nosec G201 - This is a test helper with controlled internal input, not user input
	dropSQL := fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName)
	if _, err := pgDB.Exec(dropSQL); err != nil {
		t.Logf("Warning: failed to drop test database %s: %v", dbName, err)
	} else {
		t.Logf("Successfully dropped test database: %s", dbName)
	}
}
