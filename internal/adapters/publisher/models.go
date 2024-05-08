package publisher

import (
	"context"
	"time"
)

type PublisherType string

const (
	PublisherTypeRabbitMQ PublisherType = "rabbitmq"
)

type PublisherConfig struct {
	Type  PublisherType
	Topic string
	Addr  string
}

type IPublisherConfigFactory interface {
	Config() PublisherConfig
}

type Event struct {
	ID         string    `json:"id"`
	Collection string    `json:"collection"`
	Action     string    `json:"action"`
	Data       any       `json:"data"`
	DataOld    any       `json:"data_old"`
	EventTime  time.Time `json:"event_time"`
}

type IEventPublisher interface {
	Publish(context.Context, string, *Event) error
	Close() error
}
