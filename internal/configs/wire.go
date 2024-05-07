package configs

import "github.com/google/wire"

// ConfigsSet contains project configs
var ConfigsSet = wire.NewSet(
	NewAppConfig,
	NewConfigurator,
)
