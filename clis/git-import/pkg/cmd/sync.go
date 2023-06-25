package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hacky-stuff/new-dashboards/clis/git-import/pkg/backend/es"
	"github.com/hacky-stuff/new-dashboards/clis/git-import/pkg/types"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Args:  cobra.ExactArgs(1),
	Short: "Sync one Git repository commits into ElasticSearch.",
	Run: func(cmd *cobra.Command, args []string) {
		es, err := es.GetTypedClient()
		if err != nil {
			log.Fatalf("Error creating the client: %s", err)
		}

		res, err := es.Core.
			Get("git-repository", args[0]).
			Do(context.Background())
		if err != nil {
			log.Fatalf("Error getting response: %v", err)
		}

		repo := types.Repository{}
		json.Unmarshal(res.Source_, &repo)

		log.Printf("repo: %s\n", repo.URL)
	},
}
