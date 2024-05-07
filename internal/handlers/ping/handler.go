package ping_handler

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"math/rand"
)

type PingCheckHandler struct {
}

func NewPingHandler() *PingCheckHandler {
	return &PingCheckHandler{}
}

// Ping godoc
//
//	@Summary		Dummy ping
//	@Description	Just random return 200 OK or error
//	@Tags			Test Controller
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/api/v1/ping [get]
func (hch *PingCheckHandler) Handle(ctx fiber.Ctx) error {

	n := rand.Int()
	
	if n%2 == 0 {
		return errors.New("ooops")
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "pong",
	})
}
