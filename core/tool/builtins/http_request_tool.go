package declarations

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

const (
	defaultRequestTimeoutSeconds = 30
	// maxHTTPRequestBytes caps the readable response to keep the model
	// context bounded, matching the shell tool output cap.
	maxHTTPRequestBytes = 256 * 1024
)

type HTTPRequestToolInput struct {
	Method         string            `json:"method,omitempty" jsonschema_description:"HTTP method of the request. Any method is allowed; methods other than GET and HEAD require operator approval. Default is GET."`
	URL            string            `json:"url" jsonschema_description:"The absolute http or https request URL."`
	Headers        map[string]string `json:"headers,omitempty" jsonschema_description:"Optional request headers as name/value pairs."`
	Body           string            `json:"body,omitempty" jsonschema_description:"Optional request body. Recommended for POST, PUT, and PATCH."`
	TimeoutSeconds int               `json:"timeout_seconds,omitempty" jsonschema_description:"Request timeout in seconds. Default is 30 seconds."`
	MaxBodyBytes   int               `json:"max_body_bytes,omitempty" jsonschema_description:"Maximum response body bytes to return. Default is 262144."`
	AllowPrivate   bool              `json:"allow_private,omitempty" jsonschema_description:"Set true only when the URL intentionally targets a private network host."`
}

type HTTPRequestToolOutput struct {
	Status      int               `json:"status" jsonschema_description:"HTTP response status code."`
	StatusText  string            `json:"status_text" jsonschema_description:"HTTP response status line."`
	ContentType string            `json:"content_type" jsonschema_description:"Response Content-Type header value."`
	Headers     map[string]string `json:"headers" jsonschema_description:"Response headers, first value per name."`
	Body        string            `json:"body" jsonschema_description:"Response body. Text is truncated at the byte limit; non-UTF-8 bodies are omitted with is_binary=true."`
	IsBinary    bool              `json:"is_binary" jsonschema_description:"Whether the response body was detected as binary and omitted."`
	Bytes       int64             `json:"bytes" jsonschema_description:"Total body bytes received."`
	Truncated   bool              `json:"truncated" jsonschema_description:"Whether the body was cut at the byte limit."`
}

var HTTPRequestTool = tool.DefineTool(
	"http_request_tool",
	"Perform an arbitrary http or https request and return the status, headers, and body. Loopback and link-local targets are blocked unless allow_private is set; non-GET methods require approval.",
	func(ctx *context.ExecuteContext, args *HTTPRequestToolInput, bindings tool.Bindings) (*HTTPRequestToolOutput, error) {
		if args == nil {
			return nil, fmt.Errorf("arguments cannot be empty")
		}

		parsed, err := GuardWebAccess(args.URL, args.AllowPrivate)
		if err != nil {
			return nil, err
		}

		method := strings.ToUpper(strings.TrimSpace(args.Method))
		if method == "" {
			method = http.MethodGet
		}

		// Non-idempotent requests can change external state, so they always
		// pass through the same approval gate as shell commands.
		if method != http.MethodGet && method != http.MethodHead {
			if err := requireApproval(ctx, bindings, fmt.Sprintf("Allow the agent to send HTTP %s to %s", method, parsed.Redacted())); err != nil {
				return nil, err
			}
		}

		timeout := args.TimeoutSeconds
		if timeout <= 0 {
			timeout = defaultRequestTimeoutSeconds
		}
		maxBytes := int64(args.MaxBodyBytes)
		if maxBytes <= 0 || maxBytes > maxHTTPRequestBytes {
			maxBytes = maxHTTPRequestBytes
		}

		var body io.Reader
		if args.Body != "" {
			body = strings.NewReader(args.Body)
		}

		req, err := http.NewRequest(method, parsed.String(), body)
		if err != nil {
			return nil, fmt.Errorf("failed to build request: %w", err)
		}
		for name, value := range args.Headers {
			req.Header.Set(name, value)
		}

		resp, err := (&http.Client{Timeout: time.Duration(timeout) * time.Second}).Do(req)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()

		raw, truncated := readBodyCapped(resp.Body, maxBytes)
		headers := map[string]string{}
		for name, values := range resp.Header {
			if len(values) > 0 {
				headers[name] = values[0]
			}
		}

		contentType := resp.Header.Get("Content-Type")
		isBinary := !utf8.Valid(raw)
		responseBody := string(raw)
		if isBinary {
			responseBody = ""
		}

		return &HTTPRequestToolOutput{
			Status:      resp.StatusCode,
			StatusText:  resp.Status,
			ContentType: contentType,
			Headers:     headers,
			Body:        responseBody,
			IsBinary:    isBinary,
			Bytes:       int64(len(raw)),
			Truncated:   truncated,
		}, nil
	},
	[]string{tool.NetworkingCategory, tool.ResearchCategory},
)

// readBodyCapped reads the response body up to maxBytes, reporting whether
// the limit forced a cut. When a cut happens mid-rune the trailing partial
// sequence is dropped so the text path stays UTF-8 friendly; bodies that
// remain invalid are returned as-is for the caller to flag as binary.
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
		for i := 0; i < 4 && len(data) > 0; i++ {
			data = data[:len(data)-1]
			if utf8.Valid(data) {
				break
			}
		}
	}

	return data, truncated
}
