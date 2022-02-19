package server

import (
	"fmt"
	"log"

	"github.com/Alma-media/bloxroute/pkg/transport"
	"github.com/Alma-media/bloxroute/pkg/worker"
	"github.com/spf13/cobra"
)

type Listener interface {
	Listen()
}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "daemon",
		Short:        "start testing bloxroute server",
		Aliases:      []string{"d", "server", "start"},
		SilenceUsage: true,
		RunE:         run,
	}

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	fmt.Println("server started") // TODO

	consumer := transport.GetConsumerFromContext(cmd.Context())

	log.Printf("CONSUMER: %#v", consumer)

	w := worker.New(nil)

	return w.RunParallel(
		cmd.Context(),
		consumer.Consume(cmd.Context()),
		5,
	)
}
