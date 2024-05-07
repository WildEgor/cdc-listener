package error_handler

import (
	"errors"
	"github.com/gofiber/fiber/v3"
)

// ErrorsHandler acts like global error handler
type ErrorsHandler struct {
}

func NewErrorsHandler() *ErrorsHandler {
	return &ErrorsHandler{}
}

func (hch *ErrorsHandler) Handle(ctx fiber.Ctx, err error) error {
	sc := fiber.StatusInternalServerError
	var e *fiber.Error
	if errors.As(err, &e) {
		sc = e.Code
	}

	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	ctx.Status(sc)

	// TODO: replace with your own structure if needed
	return ctx.SendString(err.Error())
}
