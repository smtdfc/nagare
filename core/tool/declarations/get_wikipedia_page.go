package declarations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/core/utils"
)

type GetWikipediaPageArgs struct {
	// Title  string `json:"title" jsonschema:"description=Wikipedia page title"`
	PageID int `json:"query" jsonschema:"description=Wikipedia page ID"`
}

type WikiExtractResp struct {
	Query struct {
		Pages map[string]struct {
			PageID  int    `json:"pageid"`
			Title   string `json:"title"`
			Extract string `json:"extract"`
		} `json:"pages"`
	} `json:"query"`
}

func getWikipediaPage(pageID int) (*WikiExtractResp, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req := utils.CreateWikipediaContentEndpoint(pageID)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data WikiExtractResp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

var get_wikipedia_page = tool.DeclareTool(
	"get_wikipedia_page",
	"Retrieve detailed content or the latest information from a specific Wikipedia page.",
	func(ctx context.Context, args GetWikipediaPageArgs) (any, error) {

		data, err := getWikipediaPage(args.PageID)
		if err != nil {
			return nil, err
		}

		// appLogger.Info(fmt.Sprintf("%v", data))
		pageContent := data.Query.Pages[fmt.Sprintf("%d", args.PageID)]
		return map[string]any{
			"query": args.PageID,
			"found": true,
			"results": map[string]any{
				"title":   pageContent.Title,
				"content": pageContent.Extract,
			},
		}, nil
	},
)
