package core

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"project-2026-06-misoastory-be-go/internal/config"
	"project-2026-06-misoastory-be-go/internal/common/middleware"
	"project-2026-06-misoastory-be-go/internal/common/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func NewHTTPServer(lc fx.Lifecycle, cfg *config.Config, db *gorm.DB) *gin.Engine {
	// AutoMigrate database
	if err := db.AutoMigrate(
		&models.Position{},
		&models.User{},
		&models.Location{},
		&models.Category{},
		&models.Permission{},
		&models.PositionPermission{},
		&models.Product{},
		&models.ProductLocation{},
		&models.JobPost{},
	); err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandler())

	addr := fmt.Sprintf(":%d", cfg.Port)

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("listen: %s\n", err)
				}
			}()

			env := os.Getenv("APP_ENV")
			if env == "" {
				env = "development"
			}

			fmt.Printf("\n🚀 Application is running on: http://localhost:%d/\n", cfg.Port)
			fmt.Printf("🌍 Environment: %s\n\n", env)
			fmt.Printf("📚 Swagger Documentaion on : http://localhost:%d/swagger/index.html\n", cfg.Port)
			fmt.Printf("📊 Health check: http://localhost:%d/api/v1/health\n", cfg.Port)
			fmt.Printf("🏓 Ping endpoint: http://localhost:%d/api/v1/ping\n\n", cfg.Port)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down gracefully...")
			return srv.Shutdown(ctx)
		},
	})

	return router
}
