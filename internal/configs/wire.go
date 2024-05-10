package configs

import (
	"github.com/WildEgor/cdc-listener/internal/adapters/publisher"
	"github.com/google/wire"
)

// ConfigsSet contains project configs
var ConfigsSet = wire.NewSet(
	NewConfigurator,
	NewAppConfig,
	NewMongoConfig,
	NewListenerConfig,
	NewPublisherConfig,
	NewLoggerConfig,
	NewMetricsConfig,
	wire.Bind(new(publisher.IPublisherConfigFactory), new(*PublisherConfig)),
)
