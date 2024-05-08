package publisher

import (
	"fmt"
	"github.com/wagslane/go-rabbitmq"
	"log/slog"
)

var _ rabbitmq.Logger = (*RabbitPublisherLogger)(nil)

type RabbitPublisherLogger struct {
}

func (t RabbitPublisherLogger) Fatalf(s string, i ...interface{}) {
	slog.Error(fmt.Sprintf(s, i))
	panic(fmt.Sprintf(s, i))
}

func (t RabbitPublisherLogger) Errorf(s string, i ...interface{}) {
	slog.Error(fmt.Sprintf(s, i))
}

func (t RabbitPublisherLogger) Warnf(s string, i ...interface{}) {
	slog.Warn(fmt.Sprintf(s, i))
}

func (t RabbitPublisherLogger) Infof(s string, i ...interface{}) {
	slog.Info(fmt.Sprintf(s, i))
}

func (t RabbitPublisherLogger) Debugf(s string, i ...interface{}) {
	slog.Debug(fmt.Sprintf(s, i))
}
