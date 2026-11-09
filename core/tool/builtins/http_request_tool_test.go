package declarations

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPRequestTool_Get(t *testing.T) {
	withBindings(t, HTTPRequestTool, &fakeApprovalBindings{approve: true})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		if r.Header.Get("X-Test") != "yes" {
			t.Errorf("expected custom header to be forwarded")
		}
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "response body")
	}))
	defer server.Close()

	res, err := HTTPRequestTool.Execute(nil, fmt.Sprintf(
		`{"method":"get","url":"%s/","headers":{"X-Test":"yes"},"allow_private":true}`,
		server.URL,
	))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(res, "200") || !contains(res, "response body") {
		t.Errorf("unexpected result: %s", res)
	}
}

func TestHTTPRequestTool_PostRequiresApproval(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: true}
	withBindings(t, HTTPRequestTool, approvals)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.Body == nil {
			t.Errorf("unexpected request method: %s", r.Method)
		}
		fmt.Fprint(w, "posted")
	}))
	defer server.Close()

	res, err := HTTPRequestTool.Execute(ctx, fmt.Sprintf(
		`{"method":"post","url":"%s/submit","body":"hello","allow_private":true}`,
		server.URL,
	))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(approvals.calls) != 1 {
		t.Errorf("expected exactly one approval request for POST, got %d", len(approvals.calls))
	}
	if !contains(res, "posted") {
		t.Errorf("unexpected result: %s", res)
	}
}

func TestHTTPRequestTool_BinaryBodyFlagged(t *testing.T) {
	withBindings(t, HTTPRequestTool, &fakeApprovalBindings{approve: true})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write([]byte{0x00, 0xff, 0x00})
	}))
	defer server.Close()

	res, err := HTTPRequestTool.Execute(nil, fmt.Sprintf(
		`{"url":"%s/blob","allow_private":true}`, server.URL,
	))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(res, `"is_binary":true`) {
		t.Errorf("expected the binary body to be flagged: %s", res)
	}
}

func TestHTTPRequestTool_Guard(t *testing.T) {
	withBindings(t, HTTPRequestTool, &fakeApprovalBindings{approve: true})

	cases := []struct {
		name string
		args string
	}{
		{"loopback without allow_private", `{"url":"http://127.0.0.1:8080/x"}`},
		{"scheme not allowed", `{"method":"get","url":"ftp://example.com/x"}`},
		{"host missing", `{"method":"get","url":"http:///x"}`},
	}
	for _, tc := range cases {
		if _, err := HTTPRequestTool.Execute(nil, tc.args); err == nil {
			t.Errorf("expected error for case '%s'", tc.name)
		}
	}
}

func TestHTTPRequestTool_AcceptsPrivateWhenAllowed(t *testing.T) {
	withBindings(t, HTTPRequestTool, &fakeApprovalBindings{approve: true})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "loopback is fine here")
	}))
	defer server.Close()

	if _, err := HTTPRequestTool.Execute(nil, fmt.Sprintf(
		`{"url":"%s/","allow_private":true}`, server.URL,
	)); err != nil {
		t.Fatalf("expected loopback targets to be allowed with allow_private=true: %v", err)
	}
}
