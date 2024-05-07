package configs

import (
	"fmt"
	"github.com/spf13/viper"
	"log/slog"
)

// ListenerConfig holds the main app configurations
type ListenerConfig struct {
	Filter []struct {
		Name        string              `mapstruct:"db"`
		Collections map[string][]string `mapstruct:"collections"`
	} `mapstruct:"filter"`

	MappedFilter map[string][]string
}

func NewListenerConfig(c *Configurator) *ListenerConfig {
	cfg := ListenerConfig{
		MappedFilter: make(map[string][]string),
	}

	if err := viper.UnmarshalKey("listener", &cfg); err != nil {
		slog.Error("app listener parse error")
	}

	for _, s := range cfg.Filter {
		for coll, ops := range s.Collections {
			cfg.MappedFilter[cfg.GetSubject(s.Name, coll)] = ops
		}
	}

	return &cfg
}

func (c *ListenerConfig) GetSubject(db, coll string) string {
	return fmt.Sprintf("%s.%s", db, coll)
}
