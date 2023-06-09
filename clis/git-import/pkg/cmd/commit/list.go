package commit

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/hacky-stuff/new-dashboards/clis/git-import/pkg/backend/es"
	"github.com/hacky-stuff/new-dashboards/clis/git-import/pkg/types"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Args:  cobra.ExactArgs(0),
	Short: "List all commits are saved in ElasticSearch.",
	Run: func(cmd *cobra.Command, args []string) {
		es, err := es.GetTypedClient()
		if err != nil {
			log.Fatalf("Error creating the client: %s", err)
		}

		res, err := es.Core.Search().
			Index("git-commit").
			// Sort("name").
			Do(context.Background())
		if err != nil {
			log.Fatalf("Error getting response: %v", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{
			"ID",
			"Hash",
			"Author",
			"Authored",
			"Message (first line)",
		})

		// log.Printf("Response: %v\n", res)
		for i := 0; i < len(res.Hits.Hits); i++ {
			hit := res.Hits.Hits[i]
			id := hit.Id_
			commit := types.Commit{}
			json.Unmarshal(hit.Source_, &commit)

			t.AppendRow([]interface{}{
				id,
				commit.Hash,
				commit.Author.Name,
				commit.Author.Date,
				commit.FirstLine,
			})
		}
		t.Render()
	},
}
