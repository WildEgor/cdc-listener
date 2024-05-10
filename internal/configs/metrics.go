package configs

import (
	"github.com/spf13/viper"
	"log/slog"
)

// MetricsConfig holds the main app configurations
type MetricsConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

func NewMetricsConfig(c *Configurator) *MetricsConfig {
	cfg := MetricsConfig{}

	if err := viper.UnmarshalKey("monitoring", &cfg); err != nil {
		slog.Error("app metrics parse error")
	}

	return &cfg
}

// IsEnabled Checks of metrics enabled
func (c *MetricsConfig) IsEnabled() bool {
	return c.Enabled
}
