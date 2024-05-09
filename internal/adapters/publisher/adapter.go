package publisher

import (
	"log/slog"
	"os"
)

// EventPublisherAdapter use as wrapper
type EventPublisherAdapter struct {
	Publisher IEventPublisher
}

func NewEventPublisherAdapter(cfg IPublisherConfigFactory) *EventPublisherAdapter {
	config := cfg.Config()
	adapter := &EventPublisherAdapter{}

	if config.Type == PublisherTypeRabbitMQ {
		pub, err := NewRabbitPublisher(cfg)
		if err != nil {
			slog.Error("publisher init error", slog.Any("err", err))
			os.Exit(1)
		}

		adapter.Publisher = pub

		return adapter
	}

	return nil
}
