package main

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Git repository configuration that are saved in ElasticSearch.",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
