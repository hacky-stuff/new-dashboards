package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hacky-stuff/new-dashboards/clis/git-import/pkg/backend/es"
	"github.com/hacky-stuff/new-dashboards/clis/git-import/pkg/types"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Args:  cobra.ExactArgs(1),
	Short: "Sync one Git repository commits into ElasticSearch.",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFlags(0)

		es, err := es.GetTypedClient()
		if err != nil {
			log.Fatalf("Error creating the client: %s", err)
		}

		res, err := es.Core.
			Get("git-repository", args[0]).
			Do(context.Background())
		if err != nil {
			log.Fatalf("Error getting response: %v", err)
		}

		id := res.Id_
		repo := types.Repository{}
		json.Unmarshal(res.Source_, &repo)

		log.Printf("repo: %s\n", repo.URL)

		path := filepath.Join("git-data", id)

		var r *git.Repository

		if _, err := os.Stat(path); os.IsNotExist(err) {
			Info("git clone %s %s", repo.URL, path)
			r, err = git.PlainClone(path, true, &git.CloneOptions{
				URL: repo.URL,
			})
			CheckIfError(err)
		} else {
			Info("plain open %s", path)
			r, err = git.PlainOpen(path)
			CheckIfError(err)
		}

		commitCount := 0
		commitIter, err := r.Log(&git.LogOptions{
			// From:  commit.Hash,
			Order: git.LogOrderCommitterTime,
		})
		CheckIfError(err)

		err = commitIter.ForEach(func(c *object.Commit) error {
			commitCount++
			commitFiles, err := c.Files()
			CheckIfError(err)

			var binaryFileCount int
			var nonBinaryFileCount int
			var overallLines int
			var overallFileSize int64
			err = commitFiles.ForEach(func(f *object.File) error {
				bin, err := f.IsBinary()
				if err == nil {
					if bin {
						binaryFileCount++
					} else {
						nonBinaryFileCount++
						lines, err := f.Lines()
						if err == nil {
							overallLines += len(lines)
						}
					}
				}
				overallFileSize += f.Size
				return nil
			})
			CheckIfError(err)

			commit := types.Commit{
				Hash: c.Hash.String(),
				Author: types.Signature{
					Name:  c.Author.Name,
					Email: c.Author.Email,
					Date:  c.Author.When,
				},
				Commit: types.Signature{
					Name:  c.Committer.Name,
					Email: c.Committer.Email,
					Date:  c.Committer.When,
				},
				FirstLine:   strings.Split(c.Message, "\n")[0],
				FullMessage: c.Message,
				TreeInfo: types.TreeInfo{
					BinaryFileCount:    binaryFileCount,
					NonBinaryFileCount: nonBinaryFileCount,
					OverallLines:       overallLines,
					OverallFileSize:    overallFileSize,
				},
			}

			for i := 0; i < len(c.ParentHashes); i++ {
				commit.ParentHashs = append(commit.ParentHashs, string(c.ParentHashes[i].String()))
			}

			fileStats, err := c.Stats()
			CheckIfError(err)
			commit.ChangeSummary.FileCount += len(fileStats)
			for i := 0; i < len(fileStats); i++ {
				commit.ChangeSummary.Addition += fileStats[i].Addition
				commit.ChangeSummary.Deletion += fileStats[i].Deletion
				commit.FileChanges = append(commit.FileChanges, types.FileChange{
					Filename: fileStats[i].Name,
					Addition: fileStats[i].Addition,
					Deletion: fileStats[i].Deletion,
				})
			}

			res, err := es.Core.Index("git-commit").
				Id(c.Hash.String()).
				Request(commit).
				Do(context.Background())
			if err != nil {
				log.Fatalf("Error getting response: %v", err)
			}

			log.Printf("Response: %v\n", res)

			return nil
		})
		CheckIfError(err)

		Info("commit count %d", commitCount)
	},
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}
