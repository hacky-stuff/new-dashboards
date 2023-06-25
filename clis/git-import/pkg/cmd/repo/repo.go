package repo

import (
	"github.com/spf13/cobra"
)

var RepoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Add, list and delete git repositories.",
}

func init() {
	RepoCmd.AddCommand(addCmd)
	RepoCmd.AddCommand(listCmd)
	RepoCmd.AddCommand(deleteCmd)
}
