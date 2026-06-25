package declarations

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/core/utils/endpoints"
)

type WikiSearchArgs struct {
	Query string `json:"query" jsonschema:"description=Wikipedia search keywords"`
}

type WikiSearchResp struct {
	Query struct {
		Search []struct {
			Title   string `json:"title"`
			PageID  int    `json:"pageid"`
			Snippet string `json:"snippet"`
		} `json:"search"`
	} `json:"query"`
}

func searchWikipedia(query string) (*WikiSearchResp, error) {

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req := endpoints.CreateWikipediaSearchEndpoint(query)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data WikiSearchResp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

var search_wikipedia = tool.DeclareTool(
	"search_wikipedia",
	"Search for information on Wikipedia.",
	func(ctx context.Context, args WikiSearchArgs) (any, error) {

		data, err := searchWikipedia(args.Query)
		if err != nil {
			return nil, err
		}

		if len(data.Query.Search) == 0 {
			return map[string]any{
				"query": args.Query,
				"found": false,
			}, nil
		}

		results := make([]map[string]any, 0)

		for i, item := range data.Query.Search {
			if i >= 3 {
				break
			}

			results = append(results, map[string]any{
				"title":      item.Title,
				"snippet":    item.Snippet,
				"page_title": item.Title,
				"page_id":    item.PageID,
			})
		}

		return map[string]any{
			"query":   args.Query,
			"found":   true,
			"results": results,
		}, nil
	},
	tool.STATIC_TOOL,
	tool.NO_GROUP,
)
