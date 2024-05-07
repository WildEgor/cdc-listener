package publisher

func NewPublisher(cfg IPublisherConfigFactory) IEventPublisher {

	config := cfg.Config()

	switch config.Type {
	case PublisherTypeRabbitMQ:
		pub, err := NewRabbitPublisher(cfg)
		if err != nil {
			return nil
		}

		return pub
	default:
		return nil
	}
}
