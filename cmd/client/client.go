package client

import (
	"github.com/Alma-media/bloxroute/cmd/config"
	"github.com/Alma-media/bloxroute/pkg/model"
	"github.com/Alma-media/bloxroute/pkg/transport"
	"github.com/Alma-media/bloxroute/pkg/transport/sqs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awsSQS "github.com/aws/aws-sdk-go/service/sqs"
	"github.com/spf13/cobra"
)

var (
	key, value string
	publisher  transport.Publisher
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "client",
		Aliases: []string{"cli"},
		Short:   "start testing bloxroute client",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var (
				ctx  = cmd.Context()
				cfg  = config.GetFromContext(ctx)
				sess = session.Must(
					session.NewSessionWithOptions(
						session.Options{
							Config: aws.Config{
								Credentials: credentials.NewStaticCredentials(cfg.ID, cfg.Secret, ""),
								Region:      &cfg.Region,
								Endpoint:    &cfg.SQSEndpoint,
							},
						},
					),
				)
				sqsClient = awsSQS.New(sess)
			)

			publisher = sqs.NewPublisher(sqsClient, cfg.QueueUrl)

			return nil
		},
		SilenceUsage: true,
	}

	cmd.AddCommand(addCommand())
	cmd.AddCommand(delCommand())
	cmd.AddCommand(getCommand())
	cmd.AddCommand(allCommand())

	return cmd
}

func getCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "get",
		Aliases:      []string{"list"},
		Short:        "get all items",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return publisher.
				Publish("get", model.Payload{
					Key: key,
				})
		},
	}

	cmd.Flags().StringVarP(&key, "key", "k", "", "item key")
	cmd.MarkFlagRequired("key")

	return cmd
}

func allCommand() *cobra.Command {
	return &cobra.Command{
		Use:          "all",
		Short:        "get all items",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return publisher.
				Publish("all", model.Payload{})
		},
	}
}

func addCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "add",
		Short:        "add new item",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return publisher.
				Publish("add", model.Payload{
					Key:   key,
					Value: value,
				})
		},
	}

	cmd.Flags().StringVarP(&key, "key", "k", "", "item key")
	cmd.Flags().StringVarP(&value, "value", "v", "", "item value")
	cmd.MarkFlagRequired("key")
	cmd.MarkFlagRequired("value")

	return cmd
}

func delCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "del",
		Aliases:      []string{"delete", "remove", "rm"},
		Short:        "remove item",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return publisher.
				Publish("del", model.Payload{
					Key: key,
				})
		},
	}

	cmd.Flags().StringVarP(&key, "key", "k", "", "item key")
	cmd.MarkFlagRequired("key")

	return cmd
}
