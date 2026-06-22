package endpoints

import (
	"net/http"
	"net/url"

	"github.com/smtdfc/nagare/core/utils"
)

func CreateGithubSearchEndpoint(query string) *http.Request {
	endpoint := "https://api.github.com/search/repositories"

	params := url.Values{}
	params.Set("sort", "stars")
	params.Set("order", "desc")
	if query != "" {
		params.Set("q", query)
	}

	fullURL := endpoint + "?" + params.Encode()
	req, _ := http.NewRequest("GET", fullURL, nil)

	req.Header.Set("User-Agent", utils.NAGARE_USER_AGENT)
	req.Header.Set("Accept", "application/json")

	return req
}
