package main

import (
	"context"
	"log"

	"github.com/hacky-stuff/new-dashboards/clis/git-import/pkg/backend/es"
	"github.com/hacky-stuff/new-dashboards/clis/git-import/pkg/types"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [repository-name] [repository-url]",
	Args:  cobra.ExactArgs(2),
	Short: "Add a new Git repository configuration to ElasticSearch.",
	Run: func(cmd *cobra.Command, args []string) {
		repo := types.Repository{
			Name: args[0],
			URL:  args[1],
		}

		es, err := es.GetTypedClient()
		if err != nil {
			log.Fatalf("Error creating the client: %s", err)
		}

		res, err := es.Core.Index("git-repository").
			Id(args[0]).
			Request(repo).
			Do(context.Background())
		if err != nil {
			log.Fatalf("Error getting response: %v", err)
		}

		log.Printf("Response: %v\n", res)
	},
}
