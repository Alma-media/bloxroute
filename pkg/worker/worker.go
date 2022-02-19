package worker

import (
	"context"

	"github.com/Alma-media/bloxroute/pkg/transport"
)

type HandlerFunc func(payload []byte) error

type Logger interface {
	Errorf(string, ...interface{})
	Debugf(string, ...interface{})
}

type Worker struct {
	commands map[string]HandlerFunc
	logger   Logger
}

func New(logger Logger) *Worker {
	return &Worker{
		commands: make(map[string]HandlerFunc),
		logger:   logger,
	}
}

// Handle adds a command to the worker (chainable).
func (w *Worker) Handle(cmd string, fn HandlerFunc) *Worker {
	w.commands[cmd] = fn

	return w
}

// Run the worker listening to the data input channel.
func (w *Worker) Run(ctx context.Context, input <-chan transport.Message) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-input:
			if !ok {
				return nil
			}

			handler, ok := w.commands[string(msg.Command())]
			if !ok {
				w.logger.Errorf("unknown command %q", msg.Command())

				continue
			}

			if err := handler(msg.Payload()); err != nil {
				w.logger.Errorf("%q error: %w", msg.Command(), err)

				continue
			}

			if err := msg.Consumed(); err != nil {
				w.logger.Errorf("filed to consume the message: %w", err)
			}
		}
	}
}

// RunParallel handles messages in parallel goroutines (it can spawn 1..poolSize concurrent routines).
// The advantage of current approach: we do not spawn idle goroutines waiting for messages.
func (w *Worker) RunParallel(ctx context.Context, input <-chan transport.Message, poolSize int) error {
	workers := make(chan struct{}, poolSize)

	for msg := range input {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case workers <- struct{}{}:
		}

		go func(msg transport.Message) {
			handler, ok := w.commands[string(msg.Command())]
			if !ok {
				w.logger.Errorf("unknown command %q", msg.Command())

				return
			}

			if err := handler(msg.Payload()); err != nil {
				w.logger.Errorf("%q error: %w", msg.Command(), err)

				return
			}

			if err := msg.Consumed(); err != nil {
				w.logger.Errorf("filed to consume the message: %w", err)
			}
		}(msg)
	}

	return nil
}
