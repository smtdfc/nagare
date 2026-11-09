package declarations

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadFileTool_Downloads(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: true}
	withBindings(t, DownloadFileTool, approvals)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "download payload")
	}))
	defer server.Close()

	dest := t.TempDir()
	res, err := DownloadFileTool.Execute(ctx, fmt.Sprintf(
		`{"url":"%s/hello.txt","destination":"%s","confirm":true,"allow_private":true}`,
		server.URL, dest,
	))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(approvals.calls) != 1 {
		t.Errorf("expected one approval request, got %d", len(approvals.calls))
	}

	file := filepath.Join(dest, "hello.txt")
	data, readErr := os.ReadFile(file)
	if readErr != nil || string(data) != "download payload" {
		t.Errorf("expected downloaded file content, got '%s' (err: %v)", string(data), readErr)
	}
	if !contains(res, "download payload") && !contains(res, "hello.txt") {
		t.Errorf("unexpected result: %s", res)
	}
}

func TestDownloadFileTool_ChecksumVerification(t *testing.T) {
	ctx := approvalTestContext()
	withBindings(t, DownloadFileTool, &fakeApprovalBindings{approve: true})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "checksum me")
	}))
	defer server.Close()

	dest := t.TempDir()
	sum := sha256.Sum256([]byte("checksum me"))

	// good digest
	_, err := DownloadFileTool.Execute(ctx, fmt.Sprintf(
		`{"url":"%s/f.txt","destination":"%s","filename":"good.txt","expected_sha256":"%s","confirm":true,"allow_private":true}`,
		server.URL, dest, hex.EncodeToString(sum[:]),
	))
	if err != nil {
		t.Fatalf("unexpected error on valid checksum: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "good.txt")); err != nil {
		t.Errorf("expected the verified file to be kept: %v", err)
	}

	// bad digest
	_, err = DownloadFileTool.Execute(ctx, fmt.Sprintf(
		`{"url":"%s/f.txt","destination":"%s","filename":"bad.txt","expected_sha256":"deadbeef","confirm":true,"allow_private":true}`,
		server.URL, dest,
	))
	if err == nil {
		t.Fatalf("expected checksum mismatch to fail the download")
	}
	if _, err := os.Stat(filepath.Join(dest, "bad.txt")); !os.IsNotExist(err) {
		t.Errorf("expected the mismatched file to be discarded, got err: %v", err)
	}
}

func TestDownloadFileTool_HttpErrorStatus(t *testing.T) {
	ctx := approvalTestContext()
	withBindings(t, DownloadFileTool, &fakeApprovalBindings{approve: true})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	if _, err := DownloadFileTool.Execute(ctx, fmt.Sprintf(
		`{"url":"%s/missing.bin","destination":"%s","confirm":true,"allow_private":true}`, server.URL, t.TempDir(),
	)); err == nil {
		t.Errorf("expected error for a non-2xx response")
	}
}

func TestDownloadFileTool_Validation(t *testing.T) {
	ctx := approvalTestContext()
	withBindings(t, DownloadFileTool, &fakeApprovalBindings{approve: true})

	cases := []struct {
		name string
		args string
	}{
		{"missing confirm", `{"url":"https://example.com/a.zip"}`},
		{"loopback without allow_private", `{"url":"http://127.0.0.1:1/x.zip","confirm":true}`},
		{"nested path escape", `{"url":"http://localhost/x/show?file=..","confirm":true,"allow_private":true}`},
	}
	for _, tc := range cases {
		if _, err := DownloadFileTool.Execute(ctx, tc.args); err == nil {
			t.Errorf("expected error for case '%s'", tc.name)
		}
	}
}

func TestDownloadFileTool_Denied(t *testing.T) {
	ctx := approvalTestContext()
	approvals := &fakeApprovalBindings{approve: false}
	withBindings(t, DownloadFileTool, approvals)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "nope")
	}))
	defer server.Close()

	dest := t.TempDir()
	if _, err := DownloadFileTool.Execute(ctx, fmt.Sprintf(
		`{"url":"%s/denied.txt","destination":"%s","confirm":true,"allow_private":true}`, server.URL, dest,
	)); err == nil {
		t.Fatalf("expected error when the user denies the download")
	}
	entries, err := os.ReadDir(dest)
	if err != nil || len(entries) != 0 {
		t.Errorf("expected the destination to stay empty after denial, got err: %v", err)
	}
}
