package publisher

type EventPublisherAdapter struct {
	Publisher IEventPublisher
}

func NewEventPublisherAdapter(cfg IPublisherConfigFactory) *EventPublisherAdapter {
	config := cfg.Config()
	adapter := &EventPublisherAdapter{}

	if config.Type == PublisherTypeRabbitMQ {
		pub, err := NewRabbitPublisher(cfg)
		if err != nil {
			return nil
		}

		adapter.Publisher = pub

		return adapter
	}

	return nil
}
