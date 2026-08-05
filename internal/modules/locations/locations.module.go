package locations

import "go.uber.org/fx"

var Module = fx.Module("locations",
	fx.Provide(
		NewLocationService,
		NewLocationHandler,
	),
)
