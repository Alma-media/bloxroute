package config

type Config struct {
	SQSEndpoint string `env:"SQS_ENDPOINT"`
	QueueUrl    string `env:"QUEUE_URL"`
	Region      string `env:"AWS_REGION" envDefault:"eu-central-1"`
}
