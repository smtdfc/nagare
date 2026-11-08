package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/discovery"
)

const WIKIPEDIA_BASE_URL = "https://en.wikipedia.org/w/rest.php/v1"

type WikipediaSearchResponse struct {
	Pages []WikipediaSearchResult `json:"pages"`
}

type WikipediaSearchResult struct {
	ID           int                 `json:"id"`
	Key          string              `json:"key"`
	Title        string              `json:"title"`
	Excerpt      string              `json:"excerpt"`
	MatchedTitle *string             `json:"matched_title"`
	Anchor       *string             `json:"anchor"`
	Description  string              `json:"description"`
	Thumbnail    *WikipediaThumbnail `json:"thumbnail"`
}

type WikipediaThumbnail struct {
	MimeType string `json:"mimetype"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Duration *int   `json:"duration"`
	URL      string `json:"url"`
}

type WikipediaPageExtractResponse struct {
	Query struct {
		Pages map[string]WikipediaPageExtract `json:"pages"`
	} `json:"query"`
}

type WikipediaPageExtract struct {
	PageID  int    `json:"pageid"`
	Title   string `json:"title"`
	Extract string `json:"extract"`
}

type WikipediaProvider struct {
	HTTPClient *http.Client
}

// Search implements [discovery.SourceProviderAdapter].
func (w *WikipediaProvider) Search(ctx context.Context, query string, limit int) ([]discovery.Source, error) {
	endpoint := WIKIPEDIA_BASE_URL + "/search/page"

	params := url.Values{}
	params.Set("q", query)
	params.Set("limit", fmt.Sprintf("%d", limit))

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint+"?"+params.Encode(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"User-Agent",
		"Nagare/1.0 (https://github.com/smtdfc/nagare)",
	)
	req.Header.Set("Accept", "application/json")

	resp, err := w.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"wikipedia API returned HTTP %d",
			resp.StatusCode,
		)
	}

	var result WikipediaSearchResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	var output = make([]discovery.Source, 0, len(result.Pages))
	for _, page := range result.Pages {
		output = append(output, discovery.Source{
			Provider:    "Wikipedia",
			ID:          fmt.Sprint(page.ID),
			Title:       page.Title,
			Description: page.Description,
			Url: fmt.Sprintf(
				"https://en.wikipedia.org/wiki/%s",
				url.PathEscape(page.Key),
			),
		})
	}
	return output, nil
}

func (w *WikipediaProvider) Fetch(
	ctx context.Context,
	id string,
) (*discovery.Detail, error) {
	endpoint := "https://en.wikipedia.org/w/api.php"

	params := url.Values{}
	params.Set("action", "query")
	params.Set("pageids", id)
	params.Set("prop", "extracts")
	params.Set("explaintext", "1")
	params.Set("format", "json")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint+"?"+params.Encode(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"User-Agent",
		"Nagare/1.0 (https://github.com/smtdfc/nagare)",
	)
	req.Header.Set("Accept", "application/json")

	resp, err := w.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"wikipedia API returned HTTP %d",
			resp.StatusCode,
		)
	}

	var result WikipediaPageExtractResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	page, ok := result.Query.Pages[id]
	if !ok {
		return nil, custom_errors.ErrResourceNotFound
	}

	return &discovery.Detail{
		SourceID: id,
		Content:  page.Extract,
	}, nil
}

func NewWikipediaProvider() *WikipediaProvider {
	return &WikipediaProvider{
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}
