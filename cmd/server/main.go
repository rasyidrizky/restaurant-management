package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"project-2026-06-misoastory-be-go/internal/config"
	"project-2026-06-misoastory-be-go/internal/database"
	"project-2026-06-misoastory-be-go/internal/handlers"
	"project-2026-06-misoastory-be-go/internal/middleware"
	"project-2026-06-misoastory-be-go/internal/models"
	"project-2026-06-misoastory-be-go/internal/routes"
	"project-2026-06-misoastory-be-go/internal/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// @title Misoastory API
// @version 1.0
// @description This is the backend for Misoastory built with Go, Gin, and GORM.
// @host localhost:3000
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func NewHTTPServer(lc fx.Lifecycle, cfg *config.Config, db *gorm.DB) *gin.Engine {
	// AutoMigrate database
	if err := db.AutoMigrate(
		&models.Position{},
		&models.User{},
		&models.Location{},
		&models.Category{},
		&models.Permission{},
		&models.PositionPermission{},
	); err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}

	router := gin.Default()
	addr := fmt.Sprintf(":%d", cfg.Port)

	// Create a standard http.Server so we can gracefully shut it down
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("Server starting on %s", addr)
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("listen: %s\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down gracefully...")
			return srv.Shutdown(ctx)
		},
	})

	return router
}

var CoreModule = fx.Module("core",
	fx.Provide(
		config.Load,
		database.NewDatabase,
		NewHTTPServer,
	),
)

var AuthModule = fx.Module("auth",
	fx.Provide(
		services.NewAuthService,
		middleware.NewAuthMiddleware,
		handlers.NewAuthHandler,
	),
)

var CategoryModule = fx.Module("categories",
	fx.Provide(
		services.NewCategoryService,
		handlers.NewCategoryHandler,
	),
)

var LocationModule = fx.Module("locations",
	fx.Provide(
		services.NewLocationService,
		handlers.NewLocationHandler,
	),
)

var UserModule = fx.Module("users",
	fx.Provide(
		services.NewUserService,
		handlers.NewUserHandler,
	),
)

func main() {
	fx.New(
		CoreModule,
		AuthModule,
		CategoryModule,
		LocationModule,
		UserModule,
		fx.Invoke(routes.SetupRoutes),
	).Run()
}
