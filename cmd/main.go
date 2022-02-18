package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/Alma-media/bloxroute/cmd/client"
	"github.com/Alma-media/bloxroute/cmd/server"
	"github.com/spf13/cobra"
)

var (
	version string
	build   string
	label   string
)

var (
	testString string
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bloxroute",
		Version: fmt.Sprintf("%s, build: %s %s", version, build, label),
		Short:   "short description",
		Long:    `long description`,
		PreRunE: setup,
	}

	cmd.AddCommand(client.New())
	cmd.AddCommand(server.New())
	cmd.DisableAutoGenTag = true

	cmd.PersistentFlags().StringVar(&testString, "stringvar", "", "test string flag")

	return cmd
}

func main() {
	err := NewCommand().Execute()
	if err != nil {
		os.Exit(1)
	}
}

func setup(cmd *cobra.Command, args []string) error {
	return errors.New("foo")
}
