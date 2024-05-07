package routers

import "github.com/gofiber/fiber/v3"

type Router interface {
	Setup(app *fiber.App)
}
