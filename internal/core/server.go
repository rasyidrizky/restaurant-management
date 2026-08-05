package core

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"project-2026-06-misoastory-be-go/internal/config"
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
	); err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	addr := fmt.Sprintf(":%d", cfg.Port)

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
