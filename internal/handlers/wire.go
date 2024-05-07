package handlers

import (
	error_handler "github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/handlers/errors"
	health_check_handler "github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/handlers/health_check"
	ping_handler "github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/handlers/ping"
	ready_check_handler "github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/handlers/ready_check"
	"github.com/google/wire"
)

// HandlersSet contains http/amqp/etc handlers (acts like facades)
var HandlersSet = wire.NewSet(
	error_handler.NewErrorsHandler,
	health_check_handler.NewHealthCheckHandler,
	ready_check_handler.NewReadyCheckHandler,
	ping_handler.NewPingHandler,
)
