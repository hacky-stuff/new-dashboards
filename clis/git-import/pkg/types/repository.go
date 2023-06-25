package types

type Repository struct {
	Name          string          `json:"name,omitempty"`
	URL           string          `json:"url,omitempty"`
	AuthorMapping []AuthorMapping `json:"authorMapping,omitempty"`
}

type AuthorMapping struct {
	Match   Author
	Replace Author
}

type Author struct {
	Name  string
	Email string
}
