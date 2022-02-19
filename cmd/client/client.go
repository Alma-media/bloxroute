package client

import (
	"github.com/spf13/cobra"
)

var (
	key, value string
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

	return cmd
}

func addCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "add",
		Short:        "add new item",
		SilenceUsage: true,
		RunE:         addItem,
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
		RunE:         removeItem,
	}

	cmd.Flags().StringVarP(&key, "key", "k", "", "item key")
	cmd.MarkFlagRequired("key")

	return cmd
}
