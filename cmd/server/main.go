package main

import (
	"fmt"

	"project-2026-06-misoastory-be-go/internal/core"
	"project-2026-06-misoastory-be-go/internal/config"
	"project-2026-06-misoastory-be-go/internal/common/middleware"
	"project-2026-06-misoastory-be-go/internal/modules/auth"
	"project-2026-06-misoastory-be-go/internal/modules/categories"
	"project-2026-06-misoastory-be-go/internal/modules/locations"
	"project-2026-06-misoastory-be-go/internal/modules/products"
	"project-2026-06-misoastory-be-go/internal/modules/users"
	"project-2026-06-misoastory-be-go/internal/modules/health"
	"project-2026-06-misoastory-be-go/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
)

// @title Misoastory API
// @version 1.0
// @description This is the backend for Misoastory built with Go, Gin, and GORM.
// @host localhost:3000
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func registerAllRoutes(
	r *gin.Engine,
	cfg *config.Config,
	m *middleware.AuthMiddleware,
	authH *auth.AuthHandler,
	catH *categories.CategoryHandler,
	locH *locations.LocationHandler,
	userH *users.UserHandler,
	prodH *products.ProductHandler,
) {
	// Swagger route
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", cfg.Port)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 route group
	v1 := r.Group("/api/v1")
	v1.GET("/health", health.HealthCheck)
	v1.GET("/ping", health.Ping)
	
	// Let each module register its own routes
	authH.RegisterRoutes(v1)
	catH.RegisterRoutes(v1, m)
	locH.RegisterRoutes(v1, m)
	userH.RegisterRoutes(v1, m)
	prodH.RegisterRoutes(v1, m)
}

func main() {
	fx.New(
		core.Module,
		auth.Module,
		categories.Module,
		locations.Module,
		users.Module,
		products.Module,
		health.Module,
		fx.Invoke(registerAllRoutes),
	).Run()
}
