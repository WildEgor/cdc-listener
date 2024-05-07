package routers

import (
	"github.com/gofiber/fiber/v3"
)

type SwaggerRouter struct {
}

func NewSwaggerRouter() *SwaggerRouter {
	return &SwaggerRouter{}
}

// Setup func for describe group of API Docs routes.
func (sr *SwaggerRouter) Setup(app *fiber.App) {
	// TODO: fiber v3 not impl swagger middleware now

	// *WIP*
	//app.Get("/docs/swagger.json", func(ctx fiber.Ctx) error {
	//	fasthttpadaptor.NewFastHTTPHandlerFunc(httpSwagger.Handler(
	//		httpSwagger.URL("docs/swagger.json"),
	//	))(ctx.Context())
	//	return nil
	//})
}
