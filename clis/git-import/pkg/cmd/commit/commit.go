package commit

import (
	"github.com/spf13/cobra"
)

var CommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "List git commits.",
}

func init() {
	CommitCmd.AddCommand(listCmd)
}
