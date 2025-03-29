package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// the default database configuration, same with in the docker-compose file
var defaultConfig = Config{
	Host:     "localhost",
	Port:     "5433",
	User:     "postgres",
	Password: "password",
	Name:     "cisab",
}

// GetConnection returns a database connection using environment variables or default config
func GetConnection() (*sql.DB, error) {
	config := Config{
		Host:     getEnv("DB_HOST", defaultConfig.Host),
		Port:     getEnv("DB_PORT", defaultConfig.Port),
		User:     getEnv("DB_USER", defaultConfig.User),
		Password: getEnv("DB_PASSWORD", defaultConfig.Password),
		Name:     getEnv("DB_NAME", defaultConfig.Name),
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Name,
	)

	return sql.Open("postgres", connStr)
}

// RunMigrations sets up a database connection using the provided config and runs migrations
func RunMigrations(ctx context.Context, config Config) error {
	// Connect to the database
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Name,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return errors.Wrap(err, "failed to connect to database")
	}
	defer db.Close()

	// Run migrations
	if err := Migrate(ctx, db); err != nil {
		return errors.Wrap(err, "migration failed")
	}

	return nil
}

// Migrate executes all SQL migration files in the migrations directory
// against the provided database connection. It assumes that the database is already created.
func Migrate(ctx context.Context, db *sql.DB) error {
	// Test connection
	if err := db.PingContext(ctx); err != nil {
		return errors.Wrap(err, "failed to ping database")
	}

	// Get all migration files
	entries, err := migrations.ReadDir("migrations")
	if err != nil {
		return errors.Wrap(err, "failed to read migrations directory")
	}

	// Apply migrations in a transaction
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			migrationName := entry.Name()

			// Read migration file
			migrationContent, err := migrations.ReadFile(filepath.Join("migrations", migrationName))
			if err != nil {
				return errors.Wrapf(err, "failed to read migration file %s", migrationName)
			}

			// Begin transaction
			tx, err := db.BeginTx(ctx, nil)
			if err != nil {
				return errors.Wrapf(err, "failed to begin transaction for migration %s", migrationName)
			}

			// Execute migration
			_, err = tx.ExecContext(ctx, string(migrationContent))
			if err != nil {
				if err := tx.Rollback(); err != nil {
					return errors.Wrapf(err, "failed to rollback transaction for migration %s", migrationName)
				}
				return errors.Wrapf(err, "failed to execute migration %s", migrationName)
			}

			// Commit transaction
			if err := tx.Commit(); err != nil {
				return errors.Wrapf(err, "failed to commit transaction for migration %s", migrationName)
			}

			log.Printf("Applied migration: %s", migrationName)
		}
	}

	return nil
}

// SetupTestDB creates a test database, runs migrations, and returns a connection.
// It assumes that the default database is already created and accessible.
func SetupTestDB(ctx context.Context, dbName string) (*sql.DB, error) {
	// 1. Get default config for database connection
	config := defaultConfig
	// 2. Connect to default postgres database
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.Name)
	pgdb, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to postgres")
	}
	defer pgdb.Close()

	// 3. Create test database
	_, err = pgdb.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create test database")
	}

	// 4. Connect to the new test database
	testConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, dbName)

	testDB, err := sql.Open("postgres", testConnStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to test database")
	}

	// Run migrations on the test database
	if err := Migrate(ctx, testDB); err != nil {
		testDB.Close()
		return nil, errors.Wrap(err, "failed to run migrations on test database")
	}

	return testDB, nil
}

// getEnv returns the value of the environment variable or a default value if not set
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
