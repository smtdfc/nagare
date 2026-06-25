package declarations

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/core/utils/endpoints"
)

type SearchGithubArgs struct {
	Query string `json:"query" jsonschema:"description=Github search keywords"`
}

type GithubSearchResp struct {
	Items []struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		FullName    string `json:"full_name"`
		Description string `json:"description"`
		CreatedAt   string `json:"created_at"`
		UpdateAt    string `json:"update_at"`
		PushedAt    string `json:"pushed_at"`
		Language    string `json:"language"`
		Forks       int    `json:"forks"`
		OpenIssues  int    `json:"open_issues"`
	} `json:"items"`
}

func searchGithub(query string) (*GithubSearchResp, error) {

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req := endpoints.CreateGithubSearchEndpoint(query)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data GithubSearchResp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

var search_github = tool.DeclareTool(
	"search_github",
	"Search for repository information on Github.",
	func(ctx domains.AgentContext, args SearchGithubArgs) (any, error) {

		data, err := searchGithub(args.Query)
		if err != nil {
			return nil, err
		}

		if len(data.Items) == 0 {
			return map[string]any{
				"query": args.Query,
				"found": false,
			}, nil
		}

		results := make([]any, 0)

		for i, item := range data.Items {
			if i >= 10 {
				break
			}

			results = append(results, item)
		}

		return map[string]any{
			"query":   args.Query,
			"found":   true,
			"results": results,
		}, nil
	},
	domains.STATIC_TOOL,
	domains.NO_GROUP,
)
