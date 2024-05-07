package pkg

import (
	"context"
	"fmt"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/configs"
	eh "github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/handlers/errors"
	nfm "github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/middlewares/not_found"
	"github.com/WildEgor/e-shop-fiber-microservice-boilerplate/internal/routers"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"github.com/google/wire"
	"log/slog"
	"os"
	"time"
)

var AppSet = wire.NewSet(
	NewApp,
	configs.ConfigsSet,
	routers.RouterSet,
)

// Server represents the main server configuration.
type Server struct {
	App       *fiber.App
	AppConfig *configs.AppConfig
}

func (srv *Server) Run(ctx context.Context) {
	slog.Info("server is listening")

	if err := srv.App.Listen(fmt.Sprintf(":%s", srv.AppConfig.Port), fiber.ListenConfig{
		DisableStartupMessage: false,
		EnablePrintRoutes:     false,
		OnShutdownSuccess: func() {
			slog.Debug("success shutdown service")
		},
	}); err != nil {
		slog.Error("unable to start server")
	}
}

func (srv *Server) Shutdown(ctx context.Context) {
	slog.Info("shutdown service")

	if err := srv.App.Shutdown(); err != nil {
		slog.Error("unable to shutdown server")
	}
}

func NewApp(
	ac *configs.AppConfig,
	eh *eh.ErrorsHandler,
	prr *routers.PrivateRouter,
	pbr *routers.PublicRouter,
	sr *routers.SwaggerRouter,
) *Server {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	if ac.IsProduction() {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}
	slog.SetDefault(logger)

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

	prr.Setup(app)
	pbr.Setup(app)
	sr.Setup(app)

	// 404 handler
	app.Use(nfm.NewNotFound())

	return &Server{
		App:       app,
		AppConfig: ac,
	}
}
