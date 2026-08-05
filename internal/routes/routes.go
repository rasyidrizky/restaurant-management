package routes

import (
	"project-2026-06-misoastory-be-go/internal/handlers"
	"project-2026-06-misoastory-be-go/internal/middleware"
	_ "project-2026-06-misoastory-be-go/docs" // Import swagger docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	locationHandler *handlers.LocationHandler,
	categoryHandler *handlers.CategoryHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health routes
	router.GET("/health", handlers.HealthCheck)
	
	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
		
		// Users routes
		users := v1.Group("/users")
		{
			users.GET("", authMiddleware.RequireAuth(), authMiddleware.RequirePermission("USER", "VIEW"), userHandler.GetUsers)
		}

		// Locations routes
		locations := v1.Group("/locations")
		{
			// Public read
			locations.GET("", locationHandler.GetLocations)
			locations.GET("/:id", locationHandler.GetLocationByID)
			// Protected write
			locations.POST("", authMiddleware.RequireAuth(), authMiddleware.RequirePermission("LOCATION", "ADD"), locationHandler.CreateLocation)
			locations.PATCH("/:id", authMiddleware.RequireAuth(), authMiddleware.RequirePermission("LOCATION", "UPDATE"), locationHandler.UpdateLocation)
			locations.DELETE("/:id", authMiddleware.RequireAuth(), authMiddleware.RequirePermission("LOCATION", "DELETE"), locationHandler.DeleteLocation)
		}

		// Categories routes
		categories := v1.Group("/categories")
		{
			// Public read
			categories.GET("", categoryHandler.GetCategories)
			categories.GET("/:id", categoryHandler.GetCategoryByID)
			// Protected write
			categories.POST("", authMiddleware.RequireAuth(), authMiddleware.RequirePermission("CATEGORY", "ADD"), categoryHandler.CreateCategory)
			categories.PATCH("/:id", authMiddleware.RequireAuth(), authMiddleware.RequirePermission("CATEGORY", "UPDATE"), categoryHandler.UpdateCategory)
			categories.DELETE("/:id", authMiddleware.RequireAuth(), authMiddleware.RequirePermission("CATEGORY", "DELETE"), categoryHandler.DeleteCategory)
		}
	}
}
