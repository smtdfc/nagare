package declarations

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/paths"
	"github.com/smtdfc/nagare/pkgs/system"
)

const (
	defaultDownloadTimeoutSeconds       = 120
	defaultMaxDownloadBytes       int64 = 1 << 30 // 1 GiB
	downloadUserAgent                   = "nagare-agent/1.0"
)

type DownloadFileToolInput struct {
	URL            string `json:"url" jsonschema_description:"The http or https URL of the file to download."`
	Destination    string `json:"destination,omitempty" jsonschema_description:"Optional directory or full file path to save to. Defaults to the agent artifacts directory."`
	Filename       string `json:"filename,omitempty" jsonschema_description:"Optional target filename when destination is a directory or omitted. Defaults to the basename of the URL."`
	TimeoutSeconds int    `json:"timeout_seconds,omitempty" jsonschema_description:"Request timeout in seconds. Default is 120 seconds."`
	MaxBytes       int64  `json:"max_bytes,omitempty" jsonschema_description:"Maximum file size to download in bytes. Default is 1073741824 (1 GiB). Larger downloads are refused."`
	ExpectedSha256 string `json:"expected_sha256,omitempty" jsonschema_description:"Optional sha256 digest used to verify the downloaded file; the download is deleted on mismatch."`
	AllowPrivate   bool   `json:"allow_private,omitempty" jsonschema_description:"Set true only when the URL intentionally targets a private network host."`
	Confirm        bool   `json:"confirm" jsonschema_description:"Must be true to allow the download."`
}

type DownloadFileToolOutput struct {
	URL      string `json:"url"`
	File     string `json:"file" jsonschema_description:"Absolute path of the downloaded file."`
	Bytes    int64  `json:"bytes" jsonschema_description:"Number of bytes written."`
	SHA256   string `json:"sha256,omitempty" jsonschema_description:"Digest of the downloaded content."`
	Verified bool   `json:"verified" jsonschema_description:"Whether the optional expected_sha256 check passed."`
}

var DownloadFileTool = tool.DefineTool(
	"download_file_tool",
	"Stream a file from an http or https URL to the local filesystem with size and timeout limits, blocking loopback and link-local targets unless explicitly allowed. Confirm must be true.",
	func(ctx *context.ExecuteContext, args *DownloadFileToolInput, bindings tool.Bindings) (*DownloadFileToolOutput, error) {
		if args == nil || !args.Confirm {
			return nil, fmt.Errorf("download requires confirm=true")
		}

		parsed, err := GuardWebAccess(args.URL, args.AllowPrivate)
		if err != nil {
			return nil, err
		}

		timeout := args.TimeoutSeconds
		if timeout <= 0 {
			timeout = defaultDownloadTimeoutSeconds
		}
		if args.MaxBytes <= 0 {
			args.MaxBytes = defaultMaxDownloadBytes
		}

		target, err := downloadTarget(parsed.Path, args.Destination, args.Filename)
		if err != nil {
			return nil, err
		}
		if _, err := os.Stat(target); err == nil {
			return nil, fmt.Errorf("destination already exists: %s", target)
		}

		if err := requireApproval(ctx, bindings, fmt.Sprintf("Allow the agent to download %s into %s", parsed.Redacted(), target)); err != nil {
			return nil, err
		}

		client := &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}

		req, err := http.NewRequest(http.MethodGet, parsed.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to build request: %w", err)
		}
		req.Header.Set("User-Agent", downloadUserAgent)

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("download failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("unsuccessful response status '%s' for %s", resp.Status, parsed.Redacted())
		}

		// Stream to a temporary file first so failures and digest mismatches
		// never leave a partial artifact behind.
		tmpFile := filepath.Join(filepath.Dir(target), "."+uuid.New().String()+".download")
		defer os.Remove(tmpFile)

		hasher := sha256.New()
		out, err := os.Create(tmpFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create temporary file: %w", err)
		}

		var written int64
		buf := make([]byte, 64*1024)
		for {
			n, readErr := resp.Body.Read(buf)
			if n > 0 {
				hasher.Write(buf[:n])
				written += int64(n)
				if written > args.MaxBytes {
					out.Close()
					return nil, fmt.Errorf("download exceeds the %d byte limit of this tool", args.MaxBytes)
				}
				if _, writeErr := out.Write(buf[:n]); writeErr != nil {
					out.Close()
					return nil, fmt.Errorf("failed to write download: %w", writeErr)
				}
			}
			if readErr == io.EOF {
				break
			}
			if readErr != nil {
				out.Close()
				return nil, fmt.Errorf("failed to read download stream: %w", readErr)
			}
		}
		out.Close()

		digest := hex.EncodeToString(hasher.Sum(nil))
		if args.ExpectedSha256 != "" && !strings.EqualFold(args.ExpectedSha256, digest) {
			return nil, fmt.Errorf("checksum mismatch: expected %s, got %s (file discarded)", args.ExpectedSha256, digest)
		}

		if err := os.Rename(tmpFile, target); err != nil {
			return nil, fmt.Errorf("failed to finalize download: %w", err)
		}

		return &DownloadFileToolOutput{
			URL:      parsed.Redacted(),
			File:     target,
			Bytes:    written,
			SHA256:   digest,
			Verified: args.ExpectedSha256 != "",
		}, nil
	},
	[]string{tool.NetworkingCategory, tool.FilesystemCategory},
)

// downloadTarget resolves the final absolute path for a download based on the
// optional destination (either a file path or a directory) and filename.
func downloadTarget(urlPath string, destination string, filename string) (string, error) {
	name := sanitizeFilename(filename)
	if name == "" {
		name = sanitizeFilename(filepath.Base(urlPath))
	}
	if name == "" {
		name = "download"
	}

	dir := destination
	if dir == "" {
		dir = paths.ArtifactsDir
		if dir == "" {
			dir = os.TempDir()
		}
	}

	resolved, err := system.ResolveFilePath(dir)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(resolved)
	switch {
	case err == nil && info.IsDir():
		return filepath.Join(resolved, name), nil
	case err == nil:
		// destination is an existing file path, replace it content-wise
		return resolved, nil
	default:
		if filepath.Ext(resolved) != "" {
			// treat a non-existing path with an extension as a file
			if err := os.MkdirAll(filepath.Dir(resolved), 0o755); err != nil {
				return "", fmt.Errorf("failed to create destination directory: %w", err)
			}
			return resolved, nil
		}
		if err := os.MkdirAll(resolved, 0o755); err != nil {
			return "", fmt.Errorf("failed to create destination directory: %w", err)
		}
		return filepath.Join(resolved, name), nil
	}
}

// sanitizeFilename strips path separators from a download filename to keep
// the artifact inside its directory.
func sanitizeFilename(filename string) string {
	name := filepath.Base(strings.TrimSpace(filename))
	if name == "." || name == "/" || name == string(os.PathSeparator) {
		return ""
	}
	return name
}
