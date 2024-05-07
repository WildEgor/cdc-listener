package routers

import (
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/handlers"
	"github.com/google/wire"
)

// RouterSet acts like "controllers" for routing http or etc.
var RouterSet = wire.NewSet(
	handlers.HandlersSet,
	NewPublicRouter,
	NewPrivateRouter,
	NewSwaggerRouter,
)
