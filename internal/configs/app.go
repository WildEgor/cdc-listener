package configs

import (
	"github.com/caarlos0/env/v7"
	"log/slog"
	"time"
)

// AppConfig holds the main app configurations
type AppConfig struct {
	Name    string `env:"APP_NAME" envDefault:"app"`
	RPCPort string `env:"APP_GRPC_PORT" envDefault:"8887"`
	Port    string `env:"APP_PORT" envDefault:"8888"`
	Mode    string `env:"APP_MODE,required"`

	StartTime time.Time `env:"START_TIME"`

	// Additional  metadata
	GoEnv   string `env:"GO_ENV" envDefault:"local"`
	Version string `env:"VERSION" envDefault:"local"`
}

func NewAppConfig(c *Configurator) *AppConfig {
	cfg := AppConfig{}

	if err := env.Parse(&cfg); err != nil {
		slog.Error("app config parse error")
	}

	slog.Info("envs", slog.Any("env", cfg))

	return &cfg
}

// IsProduction Check is application running in production mode
func (ac AppConfig) IsProduction() bool {
	return ac.Mode != "develop"
}
