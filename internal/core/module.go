package core

import (
	"project-2026-06-misoastory-be-go/internal/core/config"
	"project-2026-06-misoastory-be-go/internal/core/database"
	"project-2026-06-misoastory-be-go/internal/core/middleware"

	"go.uber.org/fx"
)

var Module = fx.Module("core",
	fx.Provide(
		config.Load,
		database.NewDatabase,
		middleware.NewAuthMiddleware,
		NewHTTPServer,
	),
)
