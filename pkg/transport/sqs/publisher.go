package sqs

import (
	"encoding/json"

	"github.com/Alma-media/bloxroute/pkg/model"
	"github.com/Alma-media/bloxroute/pkg/transport"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

var _ transport.Publisher = (*Publisher)(nil)

type Publisher struct {
	sqsClient sqsiface.SQSAPI
	queueURL  string
}

func NewPublisher(
	sqsClient sqsiface.SQSAPI,
	queueURL string,
) *Publisher {
	return &Publisher{
		sqsClient: sqsClient,
		queueURL:  queueURL,
	}
}

func (p *Publisher) Publish(command string, payload model.Payload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = p.sqsClient.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"command": {
				DataType:    aws.String("String"),
				StringValue: &command,
			},
		},
		MessageBody: aws.String(string(data)),
		QueueUrl:    &p.queueURL,
	})

	return err
}
