package server

import (
	"fmt"

	"github.com/spf13/cobra"
)

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

	return nil
}
