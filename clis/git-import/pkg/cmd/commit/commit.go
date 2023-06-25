package commit

import (
	"github.com/spf13/cobra"
)

var CommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Git commits.",
}

func init() {
	CommitCmd.AddCommand(listCmd)
	CommitCmd.AddCommand(statsCmd)
}
