package repo

import (
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

		res, err := es.Add("git-repository", args[0], repo)
		if err != nil {
			log.Fatalf("Error getting response: %v", err)
		}
		log.Printf("Response: %v\n", res)
	},
}
