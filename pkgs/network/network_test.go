package network

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGuardWebAccess(t *testing.T) {
	if _, err := GuardWebAccess("http://127.0.0.1:80/x", false); err == nil {
		t.Errorf("expected loopback dial to be refused without allow_private")
	}
	if _, err := GuardWebAccess("http://127.0.0.1:80/x", true); err != nil {
		t.Errorf("expected loopback to be allowed with allow_private: %v", err)
	}
	if _, err := GuardWebAccess("ftp://example.com", false); err == nil {
		t.Errorf("expected non-http schemes to be refused")
	}
	if _, err := GuardWebAccess("", false); err == nil {
		t.Errorf("expected an empty URL to be refused")
	}
}

func TestDo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "hello nagare")
	}))
	defer server.Close()

	resp, err := Do(RequestOptions{URL: server.URL})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != http.StatusOK || string(resp.Body) != "hello nagare" || resp.Truncated {
		t.Errorf("unexpected response: %+v", resp)
	}
	if resp.Headers["Content-Type"] != "text/plain" {
		t.Errorf("expected content type header, got %+v", resp.Headers)
	}
}

func TestDoTruncation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		out := strings.Repeat("adamente ", 4000)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(out))
	}))
	defer server.Close()

	resp, err := Do(RequestOptions{URL: server.URL, MaxBodyBytes: 1024})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Truncated {
		t.Errorf("expected the response to be truncated at %d bytes, got %d", resp.Body, len(resp.Body))
	}
	if len(resp.Body) != 1024 {
		t.Errorf("expected body cut to 1024 bytes, got %d", len(resp.Body))
	}
	if utf8Invalid(resp.Body) {
		t.Errorf("expected truncated text to stay valid UTF-8")
	}
}

func TestDownload(t *testing.T) {
	payload := []byte("download payload")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("User-Agent"); got != DefaultUserAgent {
			t.Errorf("expected the default user agent, got %q", got)
		}
		_, _ = w.Write(payload)
	}))
	defer server.Close()

	dest := filepath.Join(t.TempDir(), "file.txt")
	result, err := Download(server.URL, dest, DownloadOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Bytes != int64(len(payload)) {
		t.Errorf("expected %d bytes, got %d", len(payload), result.Bytes)
	}
	if data, statErr := os.ReadFile(dest); statErr != nil || string(data) != string(payload) {
		t.Errorf("unexpected destination content %q (err: %v)", data, statErr)
	}

	sum := sha256.Sum256(payload)
	if result.SHA256 != hex.EncodeToString(sum[:]) {
		t.Errorf("unexpected digest: %s", result.SHA256)
	}
}

func TestDownloadChecksumMismatchDiscardsFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("checksum me"))
	}))
	defer server.Close()

	dest := filepath.Join(t.TempDir(), "file.txt")
	if _, err := Download(server.URL, dest, DownloadOptions{ExpectedSha256: "deadbeef"}); err == nil {
		t.Fatalf("expected checksum mismatch to fail")
	}

	tempSiblings, readDirErr := os.ReadDir(filepath.Dir(dest))
	if readDirErr != nil {
		t.Fatal(readDirErr)
	}
	for _, entry := range tempSiblings {
		if strings.HasPrefix(entry.Name(), ".nagare-download-") {
			t.Errorf("expected the temporary download file to be removed, found %s", entry.Name())
		}
	}
	entries, _ := os.ReadDir(filepath.Dir(dest))
	if len(entries) != 0 {
		t.Errorf("expected the destination to stay empty after failure, got %+v", entries)
	}
}

func TestDownloadSizeCap(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		out := []byte(strings.Repeat("a", 4096))
		_, _ = w.Write(out)
	}))
	defer server.Close()

	dest := filepath.Join(t.TempDir(), "big.bin")
	if _, err := Download(server.URL, dest, DownloadOptions{MaxBytes: 16, Timeout: 5 * time.Second}); err == nil {
		t.Fatalf("expected the size cap refusal")
	}
}

func utf8Invalid(data []byte) bool {
	return len(data) > 0 && strings.ContainsRune(string(data), '\uFFFD')
}
