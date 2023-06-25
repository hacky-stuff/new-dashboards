package types

import "time"

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
