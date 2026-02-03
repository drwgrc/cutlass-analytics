package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cutlass_analytics/internal/api"
	"cutlass_analytics/internal/config"
	"cutlass_analytics/internal/database"
	"cutlass_analytics/internal/jobs"
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

	// Initialize and start scheduler (includes daily scraper and CSV poller)
	scheduler := jobs.NewScheduler(db)
	if err := scheduler.Start(); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}
	log.Println("Scheduler started successfully")

	// Run jobs once on server startup
	scheduler.RunOnce()

	defer func() {
		log.Println("Stopping scheduler...")
		scheduler.Stop()
	}()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Create HTTP server
	router := api.NewRouter(db)
	srv := &http.Server{
		Addr:    ":" + cfg.BackendPort,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.BackendPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	<-sigChan
	log.Println("Shutting down server...")

	// Shutdown server with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}