package client

import (
	"github.com/spf13/cobra"
)

var (
	inputFile  string
	outputFile string
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "client",
		Aliases:      []string{"cli"},
		Short:        "start testing bloxroute client",
		SilenceUsage: true,
	}

	cmd.AddCommand(addCommand())
	cmd.AddCommand(delCommand())
	cmd.AddCommand(getCommand())

	return cmd
}

func getCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "get",
		Aliases:      []string{"list"},
		Short:        "get all items",
		SilenceUsage: true,
		RunE:         getItems,
	}

	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "path to output file")
	cmd.MarkFlagRequired("output")

	return cmd
}

func addCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "add",
		Short:        "add new item",
		SilenceUsage: true,
		RunE:         addItem,
	}

	cmd.Flags().StringVarP(&inputFile, "input", "i", "", "path to input file")
	cmd.MarkFlagRequired("input")

	return cmd
}

func delCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "del",
		Aliases:      []string{"delete", "remove", "rm"},
		Short:        "remove item",
		SilenceUsage: true,
		RunE:         removeItem,
	}

	return cmd
}
