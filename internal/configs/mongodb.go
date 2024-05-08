package configs

import (
	"github.com/spf13/viper"
	"log/slog"
)

type MongoConfig struct {
	URI   string `mapstructure:"uri"`
	Debug bool   `mapstructure:"debug"`
}

func NewMongoConfig(c *Configurator) *MongoConfig {
	cfg := MongoConfig{}

	if err := viper.UnmarshalKey("database", &cfg); err != nil {
		slog.Error("mongo config error", "error", slog.AnyValue(err))
	}

	return &cfg
}
