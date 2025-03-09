package main

import (
	"log"

	"github.com/tyrese-r/go-home/internal/config"
	"github.com/tyrese-r/go-home/internal/handlers"
	"github.com/tyrese-r/go-home/internal/repository"
	"github.com/tyrese-r/go-home/internal/service"
	"github.com/tyrese-r/go-home/pkg/database"
)

func main() {
	// Load configuration
	cfg := config.New()

	// Initialize database
	db, err := database.NewSQLiteDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if clErr := db.Close(); clErr != nil {
			log.Printf("clError closing database: %v", clErr)
		}
	}()
	// Initialize repositories
	deviceRepo := repository.NewDeviceRepository(db)

	// Initialize services
	deviceService := service.NewDeviceService(deviceRepo)

	// Initialize HTTP handlers
	h := handlers.New(deviceService)

	// Start HTTP server
	err = h.StartServer(cfg.ServerAddress)
	if err != nil {
		log.Printf("Server failed: %v", err)
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Error closing database: %v", closeErr)
		}
		return
	}
}
