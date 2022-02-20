package transport

import (
	"context"

	"github.com/Alma-media/bloxroute/pkg/model"
)

type Message interface {
	Command() string
	Payload() []byte
	Consumed() error
}

type Publisher interface {
	Publish(command string, payload model.Payload) error
}

type Consumer interface {
	Consume(ctx context.Context) <-chan Message
}
