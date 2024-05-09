package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wagslane/go-rabbitmq"
	"log/slog"
)

var _ IEventPublisher = (*RabbitPublisher)(nil)

// RabbitPublisher represent event publisher for RabbitMQ.
type RabbitPublisher struct {
	pt        string              `wire:"-"`
	conn      *rabbitmq.Conn      `wire:"-"`
	publisher *rabbitmq.Publisher `wire:"-"`
}

// TODO: what if connection not established?
// NewRabbitPublisher create new RabbitPublisher instance.
func NewRabbitPublisher(cfg IPublisherConfigFactory) (*RabbitPublisher, error) {
	conn, err := rabbitmq.NewConn(cfg.Config().Addr)
	if err != nil {
		return nil, fmt.Errorf("no rabbit conn: %w", err)
	}

	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogger(&RabbitPublisherLogger{}),
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(cfg.Config().Topic),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
		rabbitmq.WithPublisherOptionsExchangeKind("topic"),
		rabbitmq.WithPublisherOptionsExchangeDurable,
	)
	if err != nil {
		return nil, fmt.Errorf("publisher: %w", err)
	}

	return &RabbitPublisher{
		cfg.Config().Topic,
		conn,
		publisher,
	}, nil
}

// Publish send events, implements IEventPublisher
func (p *RabbitPublisher) Publish(ctx context.Context, topic string, event *Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	exchange := topic
	if len(topic) == 0 {
		exchange = p.pt
	}

	slog.Debug("publish event",
		slog.Any("topic", exchange),
		slog.Any("event", event),
	)

	return p.publisher.PublishWithContext(
		ctx,
		body,
		[]string{""},
		//rabbitmq.WithPublishOptionsContentEncoding("utf-8"),
		rabbitmq.WithPublishOptionsExchange(exchange),
		rabbitmq.WithPublishOptionsContentType("application/json"),
	)
}

// Close represent finalization for RabbitMQ publisher.
func (p *RabbitPublisher) Close() error {
	if err := p.conn.Close(); err != nil {
		return fmt.Errorf("connection close: %w", err)
	}

	p.publisher.Close()

	return nil
}
