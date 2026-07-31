package main

import (
	"fmt"
	"log"

	"project-2026-06-misoastory-be-go/internal/config"
	"project-2026-06-misoastory-be-go/internal/database"
	"project-2026-06-misoastory-be-go/internal/models"
	"project-2026-06-misoastory-be-go/internal/routes"

	"github.com/gin-gonic/gin"
)

// @title Misoastory API
// @version 1.0
// @description This is the backend for Misoastory built with Go, Gin, and GORM.
// @host localhost:3000
// @BasePath /api/v1
func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	database.Connect(cfg.DBURL)
	if database.DB != nil {
		err := database.DB.AutoMigrate(
			&models.Position{},
			&models.User{},
			&models.Location{},
			&models.Category{},
		)
		if err != nil {
			log.Fatalf("Failed to auto migrate database: %v", err)
		}
	}

	// Set up router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server starting on %s", addr)
	
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
