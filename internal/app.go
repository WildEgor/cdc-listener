package pkg

import (
	"context"
	"fmt"
	"github.com/WildEgor/cdc-listener/internal/adapters/listener"
	"github.com/WildEgor/cdc-listener/internal/configs"
	eh "github.com/WildEgor/cdc-listener/internal/handlers/errors"
	nfm "github.com/WildEgor/cdc-listener/internal/middlewares/not_found"
	"github.com/WildEgor/cdc-listener/internal/repositories"
	"github.com/WildEgor/cdc-listener/internal/routers"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/google/wire"
	"log/slog"
	"time"
)

var AppSet = wire.NewSet(
	configs.ConfigsSet,
	routers.RouterSet,
	repositories.RepositoriesSet,
	NewApp,
)

// Server represents the main server configuration.
type Server struct {
	App       *fiber.App
	AppConfig *configs.AppConfig
	Listener  listener.IListener
}

func (srv *Server) Run(ctx context.Context) {
	slog.Info("server is listening")

	go func(ctx context.Context) {
		err := srv.Listener.Run(ctx)
		if err != nil {
			slog.Error("unable to start listener")
			return
		}
	}(ctx)

	if err := srv.App.Listen(fmt.Sprintf(":%s", srv.AppConfig.Port), fiber.ListenConfig{
		DisableStartupMessage: false,
		EnablePrintRoutes:     false,
		OnShutdownSuccess: func() {
			slog.Debug("success shutdown service")
		},
	}); err != nil {
		slog.Error("unable to start server")
		return
	}
}

func (srv *Server) Shutdown(ctx context.Context) {
	slog.Info("shutdown service")

	if err := srv.Listener.Stop(); err != nil {
		slog.Error("unable to shutdown listener")
	}

	if err := srv.App.Shutdown(); err != nil {
		slog.Error("unable to shutdown server")
	}
}

func NewApp(
	ac *configs.AppConfig,
	eh *eh.ErrorsHandler,
	pbr *routers.PublicRouter,
	sr *routers.SwaggerRouter,
	l listener.IListener,
) *Server {
	app := fiber.New(fiber.Config{
		AppName:      ac.Name,
		ErrorHandler: eh.Handle,
		Views:        html.New("./assets", ".html"),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	})

	app.Use(cors.New(cors.Config{
		AllowHeaders: "Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Authorization, Connection, Access-Control-Allow-Origin",
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))
	app.Use(recover.New())

	pbr.Setup(app)
	sr.Setup(app)

	// 404 handler
	app.Use(nfm.NewNotFound())

	return &Server{
		App:       app,
		Listener:  l,
		AppConfig: ac,
	}
}
