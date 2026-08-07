package careers

import "go.uber.org/fx"

// Module provides the dependencies for the careers domain
var Module = fx.Options(
	fx.Provide(
		NewCareerService,
		NewCareerHandler,
	),
)
