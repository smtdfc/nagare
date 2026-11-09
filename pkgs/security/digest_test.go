package security

import (
	"strings"
	"testing"
)

func TestParseDigestAlgorithm(t *testing.T) {
	cases := []struct {
		name     string
		expected DigestAlgorithm
	}{
		{name: "", expected: DigestSHA256},
		{name: "sha256", expected: DigestSHA256},
		{name: "sha1", expected: DigestSHA1},
		{name: "md5", expected: DigestMD5},
		{name: "  CRC32  ", expected: ""},
	}
	for _, tc := range cases {
		algorithm, err := ParseDigestAlgorithm(tc.name)
		if strings.Contains(strings.ToLower(tc.name), "crc") {
			if err == nil {
				t.Errorf("expected an error for algorithm '%s'", tc.name)
			}
			continue
		}
		if err != nil {
			t.Errorf("unexpected error for algorithm '%s': %v", tc.name, err)
		}
		if algorithm != tc.expected {
			t.Errorf("expected '%s', got '%s'", tc.expected, algorithm)
		}
	}
}

func TestDigestKnownVectors(t *testing.T) {
	cases := []struct {
		algorithm DigestAlgorithm
		expected  string
	}{
		{DigestSHA256, "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"},
		{DigestSHA1, "a9993e364706816aba3e25717850c26c9cd0d89d"},
		{DigestMD5, "900150983cd24fb0d6963f7d28e17f72"},
	}

	for _, tc := range cases {
		reader := strings.NewReader("abc")
		hexValue, size, err := tc.algorithm.Digest(reader)
		if err != nil {
			t.Fatalf("unexpected error for %s: %v", tc.algorithm, err)
		}
		if hexValue != tc.expected {
			t.Errorf("expected digest %s for %s, got %s", tc.expected, tc.algorithm, hexValue)
		}
		if size != 3 {
			t.Errorf("expected 3 hashed bytes for %s, got %d", tc.algorithm, size)
		}
	}
}
