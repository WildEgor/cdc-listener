package adapters

import (
	"github.com/WildEgor/cdc-listener/internal/adapters/listener"
	"github.com/WildEgor/cdc-listener/internal/adapters/publisher"
	"github.com/google/wire"
)

// AdaptersSet contains "adapters" to 3th party systems
var AdaptersSet = wire.NewSet(
	publisher.NewPublisher,
	// wire.Bind(new(publisher.IEventPublisher), new(*publisher.RabbitPublisher)),
	listener.NewListener,
	wire.Bind(new(listener.IListener), new(*listener.Listener)),
)
