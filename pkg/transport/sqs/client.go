package sqs

import (
	"context"

	"github.com/Alma-media/bloxroute/pkg/transport"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type Logger interface {
	Errorf(string, ...interface{})
	Debugf(string, ...interface{})
}

type Client struct {
	sqsClient sqsiface.SQSAPI
	queueURL  string
	timeout   int64
	batchSize int64
	logger    Logger
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Consume(ctx context.Context) <-chan transport.Message {
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
					c.logger.Errorf("failed to receive messages: %w", err)

					continue
				}

				for _, msg := range resp.Messages {
					output <- newJob(c.sqsClient, c.queueURL, msg)
				}
			}
		}
	}()

	return output
}
