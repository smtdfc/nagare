package declarations

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func sha256HexOf(t *testing.T, data []byte) string {
	t.Helper()
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func TestHashTool_Metadata(t *testing.T) {
	if HashTool.GetName() != "hash_tool" {
		t.Errorf("expected 'hash_tool', got '%s'", HashTool.GetName())
	}
	if len(HashTool.GetCategories()) != 1 || HashTool.GetCategories()[0] != "filesystem" {
		t.Errorf("expected FilesystemCategory, got %v", HashTool.GetCategories())
	}
}

func TestHashTool_Content(t *testing.T) {
	res, err := HashTool.Execute(nil, `{"content":"abc"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := sha256HexOf(t, []byte("abc"))
	if !contains(res, expected) {
		t.Errorf("expected sha256 of 'abc' (%s) in result: %s", expected, res)
	}
}

func TestHashTool_ContentKnownVectors(t *testing.T) {
	res, err := HashTool.Execute(nil, `{"content":"abc","algorithm":"md5"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sum := md5.Sum([]byte("abc"))
	if !contains(res, hex.EncodeToString(sum[:])) {
		t.Errorf("unexpected md5 result: %s", res)
	}

	res, err = HashTool.Execute(nil, `{"content":"abc","algorithm":"sha1"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sum1 := sha1.Sum([]byte("abc"))
	if !contains(res, hex.EncodeToString(sum1[:])) {
		t.Errorf("unexpected sha1 result: %s", res)
	}
}

func TestHashTool_File(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "digest.bin")
	payload := []byte("hello nagare")
	if err := os.WriteFile(tmp, payload, 0o644); err != nil {
		t.Fatalf("failed to write fixture: %v", err)
	}

	res, err := HashTool.Execute(nil, `{"path":"`+tmp+`"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(res, sha256HexOf(t, payload)) {
		t.Errorf("expected file hashing result in %s", res)
	}
}

func TestHashTool_RejectsInvalidInput(t *testing.T) {
	cases := []struct {
		name string
		args string
	}{
		{"no source", `{"algorithm":"sha256"}`},
		{"both sources", `{"path":"/tmp/a","content":"b"}`},
		{"unknown algorithm", `{"content":"b","algorithm":"crc32"}`},
		{"missing file", `{"path":"/nonexistent/nagare/missing.txt"}`},
	}
	for _, tc := range cases {
		if _, err := HashTool.Execute(nil, tc.args); err == nil {
			t.Errorf("expected error for case '%s'", tc.name)
		}
	}
}

func contains(haystack string, needle string) bool {
	return strings.Contains(haystack, needle)
}
