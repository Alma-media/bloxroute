package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Alma-media/bloxroute/cmd/client"
	"github.com/Alma-media/bloxroute/cmd/config"
	"github.com/Alma-media/bloxroute/cmd/server"
	"github.com/Alma-media/bloxroute/pkg/transport"
	"github.com/Alma-media/bloxroute/pkg/transport/sqs"
	"github.com/caarlos0/env/v6"
	"github.com/spf13/cobra"
)

const (
	description = `Client-Server application with the following requirements:
* multiple-threaded server;
* clients;
* External queue between the clients and server;

Clients:
* Should be configured from a command line or from a file (you decide);
* Can read data from a file or from a command line (you decide);
* Can request server to AddItem(), RemoveItem(), GetItem(), GetAllItems()
* Data is in the form of strings;

* Clients can be added / removed while not intefring to the server or other clients ;

Server:
* Has a data structure that holds the data in the memory (Ordered Map for C++);
* Server should be able to add an item, remove an item, get a single or all item from the data structure;`
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
		Long:    description,
		// PreRunE: setup,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("in pre run")

			return nil // errors.New("bar")
		},
	}

	cmd.AddCommand(client.New())
	cmd.AddCommand(server.New())
	cmd.DisableAutoGenTag = true

	cmd.PersistentFlags().StringVar(&testString, "stringvar", "", "test string flag")

	return cmd
}

func main() {
	var cfg config.Config
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	log.Println(cfg)

	consumer := sqs.NewClient()

	ctx := transport.AppendConsumerToContext(context.Background(), consumer)

	err := NewCommand().ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}

// func setup(cmd *cobra.Command, args []string) error {
// 	return errors.New("foo")
// }
