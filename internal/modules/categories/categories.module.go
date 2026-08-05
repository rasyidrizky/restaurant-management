package categories

import "go.uber.org/fx"

var Module = fx.Module("categories",
	fx.Provide(
		NewCategoryService,
		NewCategoryHandler,
	),
)
