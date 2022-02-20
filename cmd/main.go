package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Alma-media/bloxroute/cmd/client"
	"github.com/Alma-media/bloxroute/cmd/config"
	"github.com/Alma-media/bloxroute/cmd/server"
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
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

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bloxroute",
		Version: fmt.Sprintf("%s, build: %s %s", version, build, label),
		Short:   "short description",
		Long:    description,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			log.SetFormatter(&log.JSONFormatter{})
			log.SetOutput(os.Stdout)
			log.SetLevel(log.InfoLevel)

			return nil
		},
	}

	cmd.AddCommand(client.New())
	cmd.AddCommand(server.New())
	cmd.DisableAutoGenTag = true

	return cmd
}

func main() {
	var (
		cfg config.Config
		ctx = context.Background()
	)

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("cannot parse the config: %s", err)
	}

	if err := NewCommand().ExecuteContext(
		config.AppendToContext(ctx, cfg),
	); err != nil {
		os.Exit(1)
	}
}
