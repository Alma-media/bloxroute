package transport

import "context"

type consumerCtxKey struct{}

type Message interface {
	Command() string
	Payload() []byte
	Consumed() error
}

type Consumer interface {
	Consume(ctx context.Context) <-chan Message
}

func AppendConsumerToContext(parent context.Context, consumer Consumer) context.Context {
	return context.WithValue(parent, consumerCtxKey{}, consumer)
}

func GetConsumerFromContext(ctx context.Context) Consumer {
	consumer, _ := ctx.Value(consumerCtxKey{}).(Consumer)

	return consumer
}
