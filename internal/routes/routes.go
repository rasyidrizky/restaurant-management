package routes

import (
	"project-2026-06-misoastory-be-go/internal/handlers"
	"project-2026-06-misoastory-be-go/internal/middleware"
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
		// Auth routes
		authHandler := handlers.NewAuthHandler()
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
		
		// Users routes
		userHandler := handlers.NewUserHandler()
		users := v1.Group("/users")
		{
			users.GET("", middleware.RequireAuth(), middleware.RequirePermission("USER", "VIEW"), userHandler.GetUsers)
		}

		// Locations routes
		locationHandler := handlers.NewLocationHandler()
		locations := v1.Group("/locations")
		{
			// Public read
			locations.GET("", locationHandler.GetLocations)
			locations.GET("/:id", locationHandler.GetLocationByID)
			// Protected write
			locations.POST("", middleware.RequireAuth(), middleware.RequirePermission("LOCATION", "ADD"), locationHandler.CreateLocation)
			locations.PATCH("/:id", middleware.RequireAuth(), middleware.RequirePermission("LOCATION", "UPDATE"), locationHandler.UpdateLocation)
			locations.DELETE("/:id", middleware.RequireAuth(), middleware.RequirePermission("LOCATION", "DELETE"), locationHandler.DeleteLocation)
		}

		// Categories routes
		categoryHandler := handlers.NewCategoryHandler()
		categories := v1.Group("/categories")
		{
			// Public read
			categories.GET("", categoryHandler.GetCategories)
			categories.GET("/:id", categoryHandler.GetCategoryByID)
			// Protected write
			categories.POST("", middleware.RequireAuth(), middleware.RequirePermission("CATEGORY", "ADD"), categoryHandler.CreateCategory)
			categories.PATCH("/:id", middleware.RequireAuth(), middleware.RequirePermission("CATEGORY", "UPDATE"), categoryHandler.UpdateCategory)
			categories.DELETE("/:id", middleware.RequireAuth(), middleware.RequirePermission("CATEGORY", "DELETE"), categoryHandler.DeleteCategory)
		}
	}
}
