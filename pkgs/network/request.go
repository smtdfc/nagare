package network

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

const DefaultRequestTimeout = 30 * time.Second
const DefaultMaxBodyBytes = 256 * 1024

type RequestOptions struct {
	Method       string
	URL          string
	Headers      map[string]string
	Body         string
	Timeout      time.Duration
	MaxBodyBytes int64
}

type Response struct {
	Status      int
	StatusText  string
	ContentType string
	Headers     map[string]string
	Body        []byte
	Truncated   bool
}

func Do(options RequestOptions) (*Response, error) {
	method := strings.ToUpper(strings.TrimSpace(options.Method))
	if method == "" {
		method = http.MethodGet
	}

	var reader io.Reader
	if options.Body != "" {
		reader = strings.NewReader(options.Body)
	}

	request, err := http.NewRequest(method, options.URL, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	for name, value := range options.Headers {
		request.Header.Set(name, value)
	}

	timeout := options.Timeout
	if timeout <= 0 {
		timeout = DefaultRequestTimeout
	}
	maxBytes := options.MaxBodyBytes
	if maxBytes <= 0 {
		maxBytes = DefaultMaxBodyBytes
	}

	resp, err := (&http.Client{Timeout: timeout}).Do(request)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, truncated := readBodyCapped(resp.Body, maxBytes)

	headers := map[string]string{}
	for name, values := range resp.Header {
		if len(values) > 0 {
			headers[name] = values[0]
		}
	}

	return &Response{
		Status:      resp.StatusCode,
		StatusText:  resp.Status,
		ContentType: resp.Header.Get("Content-Type"),
		Headers:     headers,
		Body:        body,
		Truncated:   truncated,
	}, nil
}

func readBodyCapped(body io.Reader, maxBytes int64) ([]byte, bool) {
	data, err := io.ReadAll(io.LimitReader(body, maxBytes+1))
	if err != nil && err != io.EOF {
		return nil, false
	}

	truncated := int64(len(data)) > maxBytes
	if truncated {
		data = data[:maxBytes]
	}

	if truncated && !utf8.Valid(data) {
		for i := 0; i < 8 && len(data) > 0; i++ {
			data = data[:len(data)-1]
			if utf8.Valid(data) {
				break
			}
		}
	}

	return data, truncated
}
