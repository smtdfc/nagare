package security

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"strings"
)

// DigestAlgorithm names the supported hash algorithms by their common
// lowercase identifiers.
type DigestAlgorithm string

const (
	DigestSHA256 DigestAlgorithm = "sha256"
	DigestSHA1   DigestAlgorithm = "sha1"
	DigestMD5    DigestAlgorithm = "md5"
)

// ParseDigestAlgorithm resolves a digest algorithm by name, defaulting to
// sha256 when the name is empty.
func ParseDigestAlgorithm(name string) (DigestAlgorithm, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	switch name {
	case "", string(DigestSHA256):
		return DigestSHA256, nil
	case string(DigestSHA1):
		return DigestSHA1, nil
	case string(DigestMD5):
		return DigestMD5, nil
	default:
		return "", fmt.Errorf("unsupported algorithm '%s'; allowed: sha256, sha1, md5", name)
	}
}

// New returns a fresh hasher for the algorithm.
func (a DigestAlgorithm) New() hash.Hash {
	switch a {
	case DigestSHA1:
		return sha1.New()
	case DigestMD5:
		return md5.New()
	default:
		return sha256.New()
	}
}

// Digest streams data through the algorithm and returns the lowercase hex
// digest alongside the exact number of hashed bytes.
func (a DigestAlgorithm) Digest(reader io.Reader) (string, int64, error) {
	h := a.New()
	size, err := io.Copy(h, reader)
	if err != nil {
		return "", 0, fmt.Errorf("failed to hash stream: %w", err)
	}
	return hex.EncodeToString(h.Sum(nil)), size, nil
}
