package configs

import (
	"github.com/spf13/viper"
	"log/slog"
)

// MetricsConfig holds the main app configurations
type AppConfig struct {
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

func NewAppConfig(c *Configurator) *AppConfig {
	cfg := AppConfig{}

	if err := viper.UnmarshalKey("app", &cfg); err != nil {
		slog.Error("app config parse error")
	}

	return &cfg
}

// IsProduction Check is application running in production mode
func (ac AppConfig) IsProduction() bool {
	return ac.Mode != "develop"
}
