package config

import "context"

type configCtxKey struct{}

type Config struct {
	SQSEndpoint string `env:"SQS_ENDPOINT" envDefault:"http://localhost:4566"`
	QueueUrl    string `env:"QUEUE_URL" envDefault:"http://localhost:4566/000000000000/bloxroute"`
	Region      string `env:"AWS_REGION" envDefault:"eu-central-1"`
	ID          string `env:"AWS_ID" envDefault:"foo"`
	Secret      string `env:"AWS_SECRET" envDefault:"bar"`

	Server struct {
		MaxParallelProcesses        int   `env:"MAX_PARALLEL_PROCESSES" envDefault:"10"`
		MessageVisibilityTimeoutSec int64 `env:"MESSAGE_VISIBILITY_TIMEOUT_SEC" envDefault:"5"`
		MessageRecieveBatchSize     int64 `env:"MESSAGE_RECEIVE_BATCH_SIZE" envDefault:"10"`
	}
}

func AppendToContext(parent context.Context, config Config) context.Context {
	return context.WithValue(parent, configCtxKey{}, config)
}

func GetFromContext(ctx context.Context) Config {
	config, _ := ctx.Value(configCtxKey{}).(Config)

	return config
}
