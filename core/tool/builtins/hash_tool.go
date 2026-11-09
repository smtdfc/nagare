package declarations

import (
	"fmt"
	"os"
	"strings"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/security"
	"github.com/smtdfc/nagare/pkgs/system"
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

		algorithm, err := security.ParseDigestAlgorithm(args.Algorithm)
		if err != nil {
			return nil, err
		}

		if args.Content != "" {
			hexDigest, size, err := algorithm.Digest(strings.NewReader(args.Content))
			if err != nil {
				return nil, err
			}
			return &HashToolOutput{
				Algorithm: string(algorithm),
				Hex:       hexDigest,
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

		hexDigest, size, err := algorithm.Digest(file)
		if err != nil {
			return nil, err
		}
		return &HashToolOutput{
			Algorithm: string(algorithm),
			Hex:       hexDigest,
			Source:    path,
			Bytes:     size,
		}, nil
	},
	[]string{tool.FilesystemCategory},
)
