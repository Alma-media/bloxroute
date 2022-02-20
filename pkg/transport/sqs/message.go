package sqs

import (
	"github.com/Alma-media/bloxroute/pkg/transport"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

var _ transport.Message = (*job)(nil)

type job struct {
	sqsClient sqsiface.SQSAPI
	queueURL  string
	message   *sqs.Message
}

func newJob(sqsClient sqsiface.SQSAPI, queueURL string, message *sqs.Message) job {
	return job{
		sqsClient: sqsClient,
		queueURL:  queueURL,
		message:   message,
	}
}

func (j job) Command() string {
	cmd, ok := j.message.MessageAttributes["command"]
	if !ok || cmd.StringValue == nil {
		return ""
	}

	return *cmd.StringValue
}

func (j job) Payload() []byte {
	return []byte(*j.message.Body)
}

func (j job) Consumed() error {
	_, err := j.sqsClient.DeleteMessage(
		&sqs.DeleteMessageInput{
			QueueUrl:      &j.queueURL,
			ReceiptHandle: j.message.ReceiptHandle,
		},
	)

	return err
}
