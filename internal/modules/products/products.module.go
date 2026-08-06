package products

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewProductService,
		NewProductHandler,
	),
)
