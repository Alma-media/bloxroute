package sqs

import (
	"context"

	"github.com/Alma-media/bloxroute/pkg/transport"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

var _ transport.Consumer = (*Consumer)(nil)

type Logger interface {
	Errorf(string, ...interface{})
	Infof(string, ...interface{})
	Debugf(string, ...interface{})
}

type Consumer struct {
	sqsClient sqsiface.SQSAPI
	queueURL  string
	timeout   int64
	batchSize int64
	logger    Logger
}

func NewConsumer(
	sqsClient sqsiface.SQSAPI,
	queueURL string,
	timeout, batchSize int64,
	logger Logger,
) *Consumer {
	return &Consumer{
		sqsClient: sqsClient,
		queueURL:  queueURL,
		timeout:   timeout,
		batchSize: batchSize,
		logger:    logger,
	}
}

func (c *Consumer) Consume(ctx context.Context) <-chan transport.Message {
	c.logger.Infof("consumer started consuming batches up to %d messages", c.batchSize)

	output := make(chan transport.Message)

	go func() {
		defer close(output)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				resp, err := c.sqsClient.ReceiveMessageWithContext(
					ctx,
					&sqs.ReceiveMessageInput{
						AttributeNames: []*string{
							aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
						},
						MessageAttributeNames: []*string{
							aws.String(sqs.QueueAttributeNameAll),
						},
						QueueUrl:            aws.String(c.queueURL),
						MaxNumberOfMessages: aws.Int64(c.batchSize),
						VisibilityTimeout:   aws.Int64(c.timeout),
					},
				)
				if err != nil {
					c.logger.Errorf("failed to receive messages: %s", err)

					continue
				}

				c.logger.Debugf("%d messages received", len(resp.Messages))

				for _, msg := range resp.Messages {
					output <- newJob(c.sqsClient, c.queueURL, msg)
				}
			}
		}
	}()

	return output
}
