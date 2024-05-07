package configs

import (
	"github.com/WildEgor/cdc-listener/internal/adapters/publisher"
	"github.com/spf13/viper"
	"log/slog"
)

var _ publisher.IPublisherConfigFactory = (*PublisherConfig)(nil)

// PublisherConfig holds the main app configurations
type PublisherConfig struct {
	Topic       string                  `mapstructure:"topic"`
	TopicPrefix string                  `mapstructure:"topicPrefix"`
	Addr        string                  `mapstructure:"addr"`
	Type        publisher.PublisherType `mapstructure:"type"`
}

func NewPublisherConfig(c *Configurator) *PublisherConfig {
	cfg := PublisherConfig{}

	if err := viper.Unmarshal(&cfg); err != nil {
		slog.Error("app publisher parse error")
	}

	return &cfg
}

func (p *PublisherConfig) Config() publisher.PublisherConfig {
	return publisher.PublisherConfig{
		Topic: p.Topic,
		Addr:  p.Addr,
	}
}
