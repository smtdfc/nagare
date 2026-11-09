package declarations

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/system"
)

const (
	HashAlgorithmSHA256 = "sha256"
	HashAlgorithmSHA1   = "sha1"
	HashAlgorithmMD5    = "md5"
)

type HashToolInput struct {
	Path      string `json:"path,omitempty" jsonschema_description:"Optional path of the file to hash. Mutually exclusive with content."`
	Content   string `json:"content,omitempty" jsonschema_description:"Optional inline text to hash when no file is given. Mutually exclusive with path."`
	Algorithm string `json:"algorithm,omitempty" jsonschema_description:"Hash algorithm to use. Allowed: 'sha256' (default), 'sha1', 'md5'."`
}

type HashToolOutput struct {
	Algorithm string `json:"algorithm" jsonschema_description:"The hash algorithm that was applied."`
	Hex       string `json:"hex" jsonschema_description:"The resulting digest encoded as lowercase hex."`
	Source    string `json:"source" jsonschema_description:"Description of what was hashed: the file path or 'content'."`
	Bytes     int64  `json:"bytes" jsonschema_description:"Number of input bytes that were hashed."`
}

var HashTool = tool.DefineTool(
	"hash_tool",
	"Compute a sha256, sha1, or md5 digest of a file or inline content. Useful to verify download integrity and compare files without loading them into context.",
	func(ctx *context.ExecuteContext, args *HashToolInput, _ tool.Bindings) (*HashToolOutput, error) {
		if args == nil {
			return nil, fmt.Errorf("arguments cannot be empty")
		}
		if args.Path == "" && args.Content == "" {
			return nil, fmt.Errorf("either 'path' or 'content' is required")
		}
		if args.Path != "" && args.Content != "" {
			return nil, fmt.Errorf("provide either 'path' or 'content', not both")
		}

		algorithm, err := hashAlgorithm(args.Algorithm)
		if err != nil {
			return nil, err
		}

		if args.Content != "" {
			digest, size := hashData(algorithm.funcOf(), strings.NewReader(args.Content))
			return &HashToolOutput{
				Algorithm: algorithm.String(),
				Hex:       hex.EncodeToString(digest),
				Source:    "content",
				Bytes:     size,
			}, nil
		}

		path, err := system.ResolveFilePath(args.Path)
		if err != nil {
			return nil, err
		}
		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("file does not exist: %w", err)
		}
		if info.IsDir() {
			return nil, fmt.Errorf("path is a directory; only files can be hashed")
		}

		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		digest, size := hashData(algorithm.funcOf(), file)
		return &HashToolOutput{
			Algorithm: algorithm.String(),
			Hex:       hex.EncodeToString(digest),
			Source:    path,
			Bytes:     size,
		}, nil
	},
	[]string{tool.FilesystemCategory},
)

type hashAlgorithmKind int

const (
	hashAlgorithmKindSHA256 hashAlgorithmKind = iota
	hashAlgorithmKindSHA1
	hashAlgorithmKindMD5
)

func hashAlgorithm(name string) (hashAlgorithmKind, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	switch name {
	case "", HashAlgorithmSHA256:
		return hashAlgorithmKindSHA256, nil
	case HashAlgorithmSHA1:
		return hashAlgorithmKindSHA1, nil
	case HashAlgorithmMD5:
		return hashAlgorithmKindMD5, nil
	default:
		return 0, fmt.Errorf("unsupported algorithm '%s'; allowed: sha256, sha1, md5", name)
	}
}

func (k hashAlgorithmKind) String() string {
	switch k {
	case hashAlgorithmKindSHA1:
		return HashAlgorithmSHA1
	case hashAlgorithmKindMD5:
		return HashAlgorithmMD5
	default:
		return HashAlgorithmSHA256
	}
}

func (k hashAlgorithmKind) funcOf() func() hash.Hash {
	switch k {
	case hashAlgorithmKindSHA1:
		return sha1.New
	case hashAlgorithmKindMD5:
		return md5.New
	default:
		return sha256.New
	}
}

// hashData streams data through the hasher and returns the digest with the
// exact number of hashed bytes.
func hashData(newHash func() hash.Hash, reader io.Reader) ([]byte, int64) {
	h := newHash()
	size, _ := io.Copy(h, reader)
	return h.Sum(nil), size
}
