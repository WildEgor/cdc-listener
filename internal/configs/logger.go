package configs

import (
	"github.com/spf13/viper"
	"log/slog"
	"os"
)

type LogFormat string

var (
	LogJsonFormat   LogFormat = "json"
	LogPrettyFormat LogFormat = "pretty"
)

// MetricsConfig holds the main app configurations
type LoggerConfig struct {
	Level  slog.Level `mapstructure:"level"`
	Format LogFormat  `mapstructure:"format"`
}

func NewLoggerConfig(c *Configurator) *LoggerConfig {
	cfg := LoggerConfig{}

	if err := viper.UnmarshalKey("logger", &cfg); err != nil {
		slog.Error("app logger parse error")
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.Level,
	}))
	if cfg.Format != LogJsonFormat {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: cfg.Level,
		}))
	}
	slog.SetDefault(logger)

	return &cfg
}
