package publisher

import (
	"context"
	"time"
)

type PublisherType string

// TODO: add other publishers like kafka, nats...
const (
	PublisherTypeRabbitMQ PublisherType = "rabbitmq"
)

// PublisherConfig minimal config
type PublisherConfig struct {
	Type  PublisherType
	Topic string
	Addr  string
}

// IPublisherConfigFactory use for configurations
type IPublisherConfigFactory interface {
	Config() PublisherConfig
}

// Event base db changes event
type Event struct {
	ID         string    `json:"id"`
	Db         string    `json:"db"`
	Collection string    `json:"collection"`
	Action     string    `json:"action"`
	Data       any       `json:"data"`
	DataOld    any       `json:"data_old"`
	EventTime  time.Time `json:"event_time"`
}

// IEventPublisher for publishers
type IEventPublisher interface {
	Publish(context.Context, string, *Event) error
	Close() error
}
