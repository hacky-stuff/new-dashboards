package commit

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	esTypes "github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/calendarinterval"
	"github.com/hacky-stuff/new-dashboards/clis/git-import/pkg/backend/es"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Args:  cobra.ExactArgs(0),
	Short: "Stats for a git repository.",
	Run: func(cmd *cobra.Command, args []string) {
		es, err := es.GetTypedClient()
		if err != nil {
			log.Fatalf("Error creating the client: %s", err)
		}

		index := "git-commit"
		field := "author.date"
		interval := &calendarinterval.Year

		res, err := es.Search().
			Index(index).
			Request(
				&search.Request{
					Aggregations: map[string]esTypes.Aggregations{
						"date": {
							DateHistogram: &esTypes.DateHistogramAggregation{
								Field:            &field,
								CalendarInterval: interval,
							},
						},
					},
				},
			).
			Do(context.Background())
		if err != nil {
			log.Fatalf("Error getting response: %v", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{
			"Bucket",
			"Count",
		})

		dateAgg := res.Aggregations["date"].(*esTypes.DateHistogramAggregate)
		dateBuckets := dateAgg.Buckets.([]esTypes.DateHistogramBucket)

		for i := 0; i < len(dateBuckets); i++ {
			t.AppendRow([]interface{}{
				getBucket(interval, *dateBuckets[i].KeyAsString),
				dateBuckets[i].DocCount,
			})
		}
		t.Render()
	},
}

func getBucket(calendarinterval *calendarinterval.CalendarInterval, keyAsString string) string {
	keyAsString = strings.Trim(keyAsString, "\"")

	switch calendarinterval.Name {
	case "year":
		return (keyAsString)[0:4]
	case "quarter":
		switch (keyAsString)[5:10] {
		case "01-01":
			return (keyAsString)[0:4] + " Q1"
		case "04-01":
			return (keyAsString)[0:4] + " Q2"
		case "07-01":
			return (keyAsString)[0:4] + " Q3"
		case "10-01":
			return (keyAsString)[0:4] + " Q4"
		}
	case "month":
		return (keyAsString)[0:10]
	}
	return keyAsString
}
