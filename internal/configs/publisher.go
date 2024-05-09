package configs

import (
	"github.com/WildEgor/cdc-listener/internal/adapters/publisher"
	"github.com/spf13/viper"
	"log/slog"
)

var _ publisher.IPublisherConfigFactory = (*PublisherConfig)(nil)

// PublisherConfig holds the main app configurations
type PublisherConfig struct {
	// HINT: default topic
	Topic       string                  `mapstructure:"topic"`
	TopicPrefix string                  `mapstructure:"topicPrefix"`
	Addr        string                  `mapstructure:"uri"`
	Type        publisher.PublisherType `mapstructure:"type"`
}

func NewPublisherConfig(c *Configurator) *PublisherConfig {
	cfg := PublisherConfig{}

	if err := viper.UnmarshalKey("publisher", &cfg); err != nil {
		slog.Error("app publisher parse error")
	}

	return &cfg
}

func (p *PublisherConfig) Config() publisher.PublisherConfig {
	return publisher.PublisherConfig{
		Type:  p.Type,
		Topic: p.Topic,
		Addr:  p.Addr,
	}
}
