// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package pkg

import (
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/configs"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/handlers/errors"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/handlers/health_check"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/handlers/ping"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/handlers/ready_check"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/routers"
	"github.com/google/wire"
)

// Injectors from server.go:

func NewServer() (*Server, error) {
	configurator := configs.NewConfigurator()
	appConfig := configs.NewAppConfig(configurator)
	errorsHandler := error_handler.NewErrorsHandler()
	privateRouter := routers.NewPrivateRouter()
	healthCheckHandler := health_check_handler.NewHealthCheckHandler()
	readyCheckHandler := ready_check_handler.NewReadyCheckHandler()
	pingCheckHandler := ping_handler.NewPingHandler()
	publicRouter := routers.NewPublicRouter(healthCheckHandler, readyCheckHandler, pingCheckHandler)
	swaggerRouter := routers.NewSwaggerRouter()
	server := NewApp(appConfig, errorsHandler, privateRouter, publicRouter, swaggerRouter)
	return server, nil
}

// server.go:

var ServerSet = wire.NewSet(AppSet)
