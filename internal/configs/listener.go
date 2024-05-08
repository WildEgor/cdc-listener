package configs

import (
	"fmt"
	"github.com/spf13/viper"
	"log/slog"
	"strings"
)

type DbColl struct {
	Db   string
	Coll string
}

type FilterConfig struct {
	Db          string              `mapstructure:"db"`
	Collections map[string][]string `mapstructure:"collections"`
}

// ListenerConfig holds the main app configurations
type ListenerConfig struct {
	Filter       []FilterConfig    `mapstructure:"filter"`
	TopicsMap    map[string]string `mapstructure:"topicsMap"`
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
			cfg.MappedFilter[cfg.GetSubject(s.Db, coll)] = ops
		}
	}

	return &cfg
}

func (c *ListenerConfig) GetSubject(db, coll string) string {
	return fmt.Sprintf("%s.%s", db, coll)
}

func (c *ListenerConfig) GetDbCollBySubject(subj string) *DbColl {
	parsed := strings.Split(subj, ".")

	return &DbColl{
		Db:   parsed[0],
		Coll: parsed[1],
	}
}

func (c *ListenerConfig) GetTopic(subj string) string {
	parsed := strings.Split(subj, ".")
	result := strings.Join(parsed, "-")

	topic, ok := c.TopicsMap[result]
	if !ok {
		return ""
	}

	return topic
}
