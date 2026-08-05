package core

import (
	"project-2026-06-misoastory-be-go/internal/common/middleware"
	"project-2026-06-misoastory-be-go/internal/config"

	"go.uber.org/fx"
)

var Module = fx.Module("core",
	fx.Provide(
		config.Load,
		config.NewDatabase,
		middleware.NewAuthMiddleware,
		NewHTTPServer,
	),
)
