package discovery

import "context"

type Source struct {
	Provider    string `json:"provider"`
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

type Detail struct {
	SourceID string `json:"source_id"`
	Content  string `json:"content"`
}

type SourceProviderAdapter interface {
	Search(ctx context.Context, query string, limit int) ([]Source, error)
	Fetch(ctx context.Context, id string) (*Detail, error)
}
