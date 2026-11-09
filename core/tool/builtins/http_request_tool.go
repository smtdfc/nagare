package declarations

import (
	"fmt"
	"strings"
	"time"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/network"
)

// maxHTTPRequestBytes caps the readable response to keep the model context
// bounded, matching the shell tool output cap.
const maxHTTPRequestBytes = network.DefaultMaxBodyBytes

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
	Bytes       int               `json:"bytes" jsonschema_description:"Total body bytes received."`
	Truncated   bool              `json:"truncated" jsonschema_description:"Whether the body was cut at the byte limit."`
}

var HTTPRequestTool = tool.DefineTool(
	"http_request_tool",
	"Perform an arbitrary http or https request and return the status, headers, and body. Loopback and link-local targets are blocked unless allow_private is set; non-GET methods require approval.",
	func(ctx *context.ExecuteContext, args *HTTPRequestToolInput, bindings tool.Bindings) (*HTTPRequestToolOutput, error) {
		if args == nil {
			return nil, fmt.Errorf("arguments cannot be empty")
		}

		parsed, err := network.GuardWebAccess(args.URL, args.AllowPrivate)
		if err != nil {
			return nil, err
		}

		// Non-idempotent requests can change external state, so they always
		// pass through the same approval gate as shell commands.
		method := strings.ToUpper(strings.TrimSpace(args.Method))
		if method == "" {
			method = "GET"
		}
		if method != "GET" && method != "HEAD" {
			if err := requireApproval(ctx, bindings, fmt.Sprintf("Allow the agent to send HTTP %s to %s", method, parsed.Redacted())); err != nil {
				return nil, err
			}
		}

		response, err := network.Do(network.RequestOptions{
			Method:       method,
			URL:          parsed.String(),
			Headers:      args.Headers,
			Body:         args.Body,
			Timeout:      time.Duration(args.TimeoutSeconds) * time.Second,
			MaxBodyBytes: int64(args.MaxBodyBytes),
		})
		if err != nil {
			return nil, err
		}

		isBinary := !network.BodyIsText(response.Body)
		outputBody := ""
		if !isBinary {
			outputBody = string(response.Body)
		}

		return &HTTPRequestToolOutput{
			Status:      response.Status,
			StatusText:  response.StatusText,
			ContentType: response.ContentType,
			Headers:     response.Headers,
			Body:        outputBody,
			IsBinary:    isBinary,
			Bytes:       len(response.Body),
			Truncated:   response.Truncated,
		}, nil
	},
	[]string{tool.NetworkingCategory, tool.ResearchCategory},
)
