package configs

import (
	"github.com/spf13/viper"
	"log/slog"
)

// Configurator dummy
type Configurator struct{}

func NewConfigurator() *Configurator {
	c := &Configurator{}
	c.load()

	return c
}

// load Load env data from files config
func (c *Configurator) load() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		slog.Error("error loading config file")
		return
	}
}
