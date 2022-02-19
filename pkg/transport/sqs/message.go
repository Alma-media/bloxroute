package sqs

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

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
	return *j.message.Attributes["command"]
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
