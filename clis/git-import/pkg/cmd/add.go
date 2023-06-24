package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [name] [repository]",
	Args:  cobra.MinimumNArgs(2),
	Short: "Add a new Git repository configuration to ElasticSearch.",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		repo := args[1]
		fmt.Printf("Will add %s\n", name)
		fmt.Printf("  Repository: %s\n", repo)
	},
}
