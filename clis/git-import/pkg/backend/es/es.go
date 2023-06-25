package es

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetTypedClient() (*elasticsearch.TypedClient, error) {
	return elasticsearch.NewTypedClient(elasticsearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
}

func main() {
	log.SetFlags(0)
	/*
		log.Println("ElasticSearch client version:", elasticsearch.Version)

		// Initialize a client with the default settings.
		//
		// An `ELASTICSEARCH_URL` environment variable will be used when exported.
		//
		es, err := elasticsearch.NewDefaultClient()
		if err != nil {
			log.Fatalf("Error creating the client: %s", err)
		}

		// 2. Index documents
		//
		for i, title := range []string{"Test One", "Test Two"} {
			// Build the request body.
			data, err := json.Marshal(struct {
				Title string `json:"title"`
			}{
				Title: title,
			})
			if err != nil {
				log.Fatalf("Error marshaling document: %s", err)
			}

			// Set up the request object.
			req := esapi.IndexRequest{
				Index:      "test",
				DocumentID: strconv.Itoa(i + 1),
				Body:       bytes.NewReader(data),
				Refresh:    "true",
			}

			// Perform the request with the client.
			res, err := req.Do(context.Background(), es)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
			} else {
				// Deserialize the response into a map.
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					log.Printf("Error parsing the response body: %s", err)
				} else {
					// Print the response status and indexed document version.
					log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
				}
			}
		}

		// 3. Search for the indexed documents
		//
		// Build the request body.
		var buf bytes.Buffer
		query := map[string]interface{}{
			"query": map[string]interface{}{
				"match": map[string]interface{}{
					"title": "test",
				},
			},
		}
		if err := json.NewEncoder(&buf).Encode(query); err != nil {
			log.Fatalf("Error encoding query: %s", err)
		}

		// Perform the search request.
		res, err = es.Search(
			es.Search.WithContext(context.Background()),
			es.Search.WithIndex("test"),
			es.Search.WithBody(&buf),
			es.Search.WithTrackTotalHits(true),
			es.Search.WithPretty(),
		)
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			var e map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
				log.Fatalf("Error parsing the response body: %s", err)
			} else {
				// Print the response status and error information.
				log.Fatalf("[%s] %s: %s",
					res.Status(),
					e["error"].(map[string]interface{})["type"],
					e["error"].(map[string]interface{})["reason"],
				)
			}
		}

		var searchResponse map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&searchResponse); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}
		// Print the response status, number of results, and request duration.
		log.Printf(
			"[%s] %d hits; took: %dms",
			res.Status(),
			int(searchResponse["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
			int(searchResponse["took"].(float64)),
		)
		// Print the ID and document source for each hit.
		for _, hit := range searchResponse["hits"].(map[string]interface{})["hits"].([]interface{}) {
			log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
		}

		log.Println(strings.Repeat("=", 37))
	*/

	url := "https://github.com/jerolimov/react-showdown"
	path := "/tmp/react-showdown"

	var r *git.Repository

	if _, err := os.Stat(path); os.IsNotExist(err) {
		Info("git clone %s %s", url, path)
		r, err = git.PlainClone(path, false, &git.CloneOptions{
			URL: url,
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

		commit := Commit{
			Hash: c.Hash.String(),
			Author: Signature{
				Name:  c.Author.Name,
				Email: c.Author.Email,
				Date:  c.Author.When,
			},
			Commit: Signature{
				Name:  c.Committer.Name,
				Email: c.Committer.Email,
				Date:  c.Committer.When,
			},
			FirstLine:   strings.Split(c.Message, "\n")[0],
			FullMessage: c.Message,
			TreeInfo: TreeInfo{
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
			commit.FileChanges = append(commit.FileChanges, FileChange{
				Filename: fileStats[i].Name,
				Addition: fileStats[i].Addition,
				Deletion: fileStats[i].Deletion,
			})
		}

		commitJSON, err := json.MarshalIndent(commit, "", "  ")
		CheckIfError(err)
		if commitCount == 7 {
			fmt.Println(string(commitJSON))
		}

		return nil
	})
	CheckIfError(err)

	Info("commit count %d", commitCount)
}

type Signature struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

type Commit struct {
	ParentHashs   []string      `json:"parentHashs"`
	Hash          string        `json:"hash"`
	Author        Signature     `json:"author"`
	Commit        Signature     `json:"commit"`
	FirstLine     string        `json:"firstLine"`
	FullMessage   string        `json:"fullMessage"`
	ChangeSummary ChangeSummary `json:"changeSummary"`
	FileChanges   []FileChange  `json:"fileChanges"`
	TreeInfo      TreeInfo      `json:"treeInfo"`
}

type ChangeSummary struct {
	FileCount int `json:"fileCount"`
	Addition  int `json:"addition"`
	Deletion  int `json:"deletion"`
}

type FileChange struct {
	Filename string `json:"filename"`
	Addition int    `json:"addition"`
	Deletion int    `json:"deletion"`
}

type TreeInfo struct {
	BinaryFileCount    int   `json:"binaryFileCount"`
	NonBinaryFileCount int   `json:"nonBinaryFileCount"`
	OverallLines       int   `json:"overallLines"`
	OverallFileSize    int64 `json:"overallFileSize"`
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
