package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print multiple version numbers, incl. the ElasticSearch library and server.",
	Run: func(cmd *cobra.Command, args []string) {
		version := "0.0.0"
		if json, _ := cmd.Flags().GetBool("json"); json {
			fmt.Printf("{ \"version\": \"%s\" }\n", version)
		} else {
			fmt.Printf("Version: %s\n", version)
		}
	},
}

func init() {
	versionCmd.Flags().Bool("json", false, "Print verions as JSON output.")
}
