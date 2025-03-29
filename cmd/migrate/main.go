package main

import (
	"os"

	"github.com/llkhacquan/knovel-assignment/db"
	"github.com/llkhacquan/knovel-assignment/pkg/utils/logger"
)

func main() {
	var l = logger.NewDefault()
	// for simplicity, we use hardcoded config here, we can use env vars or flags in real world
	config := db.Config{
		Host:     "localhost",
		Port:     "5433",
		User:     "postgres",
		Password: "password",
		Name:     "knovel",
	}
	// something not done:
	// - incremental migration (we run all migrations at once)
	// - rollback migration

	// Run database migrations with the config
	if err := db.RunMigrations(config); err != nil {
		l.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}
	l.Info("Migrations completed successfully")
}
