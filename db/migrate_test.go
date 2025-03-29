package db

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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

// TestMigrate tests the migration functionality
func TestMigrate(t *testing.T) {
	// Generate a random database name
	dbName := generateRandomDBName()
	t.Logf("Using test database: %s", dbName)

	// Create a test database
	db, err := SetupTestDB(context.Background(), dbName)
	require.Nil(t, err, "failed to create test database")
	db.Close()
}
