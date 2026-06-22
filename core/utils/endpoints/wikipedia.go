package endpoints

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/smtdfc/nagare/core/utils"
)

func CreateWikipediaSearchEndpoint(query string) *http.Request {
	endpoint := "https://en.wikipedia.org/w/api.php"

	params := url.Values{}
	params.Set("action", "query")
	params.Set("format", "json")

	if query != "" {
		params.Set("list", "search")
		params.Set("srsearch", query)
	}

	fullURL := endpoint + "?" + params.Encode()

	req, _ := http.NewRequest("GET", fullURL, nil)

	req.Header.Set("User-Agent", utils.NAGARE_USER_AGENT)
	req.Header.Set("Accept", "application/json")

	return req
}

func CreateWikipediaContentEndpoint(pageID int) *http.Request {
	endpoint := "https://en.wikipedia.org/w/api.php"

	params := url.Values{}
	params.Set("action", "query")
	params.Set("format", "json")
	params.Set("prop", "extracts")
	params.Set("exintro", "true")
	params.Set("explaintext", "true")
	params.Set("pageids", fmt.Sprint(pageID))

	fullURL := endpoint + "?" + params.Encode()

	req, _ := http.NewRequest("GET", fullURL, nil)

	req.Header.Set("User-Agent", utils.NAGARE_USER_AGENT)
	req.Header.Set("Accept", "application/json")

	return req
}
