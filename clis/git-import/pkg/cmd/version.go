package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/info"
	"github.com/hacky-stuff/new-dashboards/clis/git-import/pkg/backend/es"
	"github.com/spf13/cobra"
)

type VersionInfo struct {
	Version       string `json:"version"`
	ElasticSearch struct {
		Client string         `json:"client"`
		Server *info.Response `json:"server"`
	} `json:"elasticSearch"`
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print multiple version numbers, incl. the ElasticSearch library and server.",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFlags(0)

		versionInfo := VersionInfo{
			Version: "unknown",
		}

		// ElasticSearch
		versionInfo.ElasticSearch.Client = elasticsearch.Version

		es, err := es.GetTypedClient()
		if err != nil {
			log.Fatalf("Error creating the client: %s", err)
		}

		infoRes, err := es.Info().Do(context.Background())
		if err != nil {
			log.Fatalf("Error getting response: %v", err)
		}
		versionInfo.ElasticSearch.Server = infoRes

		if jsonOutput, _ := cmd.Flags().GetBool("json"); jsonOutput {
			jsonBytes, _ := json.MarshalIndent(versionInfo, "", "  ")
			fmt.Println(string(jsonBytes))
		} else {
			fmt.Printf("Version: %s\n", versionInfo.Version)

			fmt.Println()
			fmt.Printf("ElasticSearch\n  Client library: %s\n", versionInfo.ElasticSearch.Client)
			if versionInfo.ElasticSearch.Server != nil {
				esVersion := versionInfo.ElasticSearch.Server.Version
				fmt.Printf("  Server version: %s\n", strings.Trim(esVersion.Int, "\""))
				fmt.Printf("  Lucene version: %s\n", esVersion.LuceneVersion)
			}
		}
	},
}

func init() {
	versionCmd.Flags().Bool("json", false, "Print verions as JSON output.")
}
