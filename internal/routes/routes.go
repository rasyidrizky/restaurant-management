package routes

import (
	"project-2026-06-misoastory-be-go/internal/handlers"
	_ "project-2026-06-misoastory-be-go/docs" // Import swagger docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine) {
	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health routes
	router.GET("/health", handlers.HealthCheck)
	
	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Auth routes placeholder
		auth := v1.Group("/auth")
		{
			auth.POST("/login", func(c *gin.Context) { c.JSON(200, gin.H{"message": "login"}) })
		}
		
		// Users routes
		userHandler := handlers.NewUserHandler()
		users := v1.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
		}

		// Locations routes
		locationHandler := handlers.NewLocationHandler()
		locations := v1.Group("/locations")
		{
			locations.GET("", locationHandler.GetLocations)
			locations.GET("/:id", locationHandler.GetLocationByID)
			locations.POST("", locationHandler.CreateLocation)
			locations.PATCH("/:id", locationHandler.UpdateLocation)
			locations.DELETE("/:id", locationHandler.DeleteLocation)
		}
	}
}
