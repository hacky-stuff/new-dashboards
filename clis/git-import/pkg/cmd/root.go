package main

import (
	"github.com/hacky-stuff/new-dashboards/clis/git-import/pkg/cmd/commit"
	"github.com/hacky-stuff/new-dashboards/clis/git-import/pkg/cmd/repo"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-import",
	Short: "git-import is a tool to import git commits into ElasticSearch.",
}

func init() {
	rootCmd.AddCommand(repo.RepoCmd)
	rootCmd.AddCommand(commit.CommitCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(versionCmd)
}
