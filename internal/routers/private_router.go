package routers

import (
	"github.com/gofiber/fiber/v3"
)

type PrivateRouter struct {
}

func NewPrivateRouter() *PrivateRouter {
	return &PrivateRouter{}
}

func (r *PrivateRouter) Setup(app *fiber.App) {
}
