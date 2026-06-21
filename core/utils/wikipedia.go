package utils

import (
	"fmt"
	"net/http"
	"net/url"
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

	req.Header.Set("User-Agent", NAGARE_USER_AGENT)
	req.Header.Set("Accept", "application/json")

	return req
}

func CreateWikipediaContentEndpoint(pageID int) *http.Request {
	endpoint := "https://en.wikipedia.org/w/api.php"

	params := url.Values{}
	params.Set("action", "query")
	params.Set("format", "json")

	// Sử dụng prop=extracts để lấy nội dung tóm tắt hoặc toàn bộ một cách tối ưu
	params.Set("prop", "extracts")

	// exintro=true: Chỉ lấy phần giới thiệu (rất hiệu quả để tránh vượt quá token)
	// Nếu bạn muốn lấy toàn bộ, hãy bỏ tham số exintro
	params.Set("exintro", "true")

	// explaintext=true: Trả về văn bản thuần (plain text), loại bỏ mã wiki phức tạp
	params.Set("explaintext", "true")

	params.Set("pageids", fmt.Sprint(pageID))

	fullURL := endpoint + "?" + params.Encode()

	req, _ := http.NewRequest("GET", fullURL, nil)

	req.Header.Set("User-Agent", NAGARE_USER_AGENT)
	req.Header.Set("Accept", "application/json")

	return req
}
