package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"cutlass_analytics/internal/api"
	"cutlass_analytics/internal/config"
	"cutlass_analytics/internal/database"
)

func main() {
	cfg := config.Load()

	// Connect to Database
	log.Println("Connecting to database...")
	db, err := database.Connect(cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Successfully connected to database")

	// Ensure database connection is closed on exit
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection closed")
		}
	}()

	// Run migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// TODO: Initialize repositories for background jobs

	// TODO: Start background jobs

	// TODO: Daily scraper at 1AM PST

	// TODO: CSV poller every 10 minutes

	router := api.NewRouter(db)
	log.Printf("Server starting on port %s", cfg.BackendPort)
	if err := router.Run(":" + cfg.BackendPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}