package repo

import (
	"log"

	"github.com/hacky-stuff/new-dashboards/clis/git-import/pkg/backend/es"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Args:  cobra.ExactArgs(1),
	Short: "Delete a Git repository from ElasticSearch.",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := es.Delete("git-repository", args[0])
		if err != nil {
			log.Fatalf("Error getting response: %v", err)
		}
		log.Printf("Result: %s\n", res.Result)
	},
}
