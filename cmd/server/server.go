package server

import (
	"fmt"
	"io"
	"os"

	"github.com/Alma-media/bloxroute/cmd/config"
	"github.com/Alma-media/bloxroute/pkg/adapter"
	"github.com/Alma-media/bloxroute/pkg/executer"
	"github.com/Alma-media/bloxroute/pkg/transport/sqs"
	"github.com/Alma-media/bloxroute/pkg/worker"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awsSQS "github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var outputFile string

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "daemon",
		Short:        "start testing bloxroute server",
		Aliases:      []string{"d", "server", "start"},
		SilenceUsage: true,
		RunE:         run,
	}

	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "server output")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	var executerOutput io.Writer = os.Stdout

	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("cannot create output file: %w", err)
		}

		defer file.Close()

		executerOutput = file
	}

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
		executer  = executer.New(adapter.New(), executerOutput)
	)

	return worker.
		New(log.WithField("component", "worker")).
		Handle("add", executer.Insert).
		Handle("del", executer.Delete).
		Handle("all", executer.GetAll).
		Handle("get", executer.GetOne).
		RunParallel(
			cmd.Context(),
			sqs.NewConsumer(
				sqsClient,
				cfg.QueueUrl,
				cfg.Server.MessageVisibilityTimeoutSec,
				cfg.Server.MessageRecieveBatchSize,
				log.WithField("component", "consumer"),
			).Consume(ctx),
			cfg.Server.MaxParallelProcesses,
		)
}
