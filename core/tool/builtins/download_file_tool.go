package declarations

import (
	"fmt"
	"os"
	"time"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/network"
	"github.com/smtdfc/nagare/pkgs/paths"
	"github.com/smtdfc/nagare/pkgs/system"
)

// defaultDownloadMaxBytes caps downloads issued by the tool: the helper
// default is 1 GiB; the tool keeps the same default but allows explicit ops.
const defaultDownloadMaxBytes = network.DefaultDownloadMaxBytes

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

		parsed, err := network.GuardWebAccess(args.URL, args.AllowPrivate)
		if err != nil {
			return nil, err
		}

		destination := args.Destination
		if destination == "" {
			destination = paths.ArtifactsDir
			if destination == "" {
				destination = os.TempDir()
			}
		}
		target, err := system.DeriveDownloadPath(parsed.Path, destination, args.Filename)
		if err != nil {
			return nil, err
		}
		if _, err := os.Stat(target); err == nil {
			return nil, fmt.Errorf("destination already exists: %s", target)
		}

		if err := requireApproval(ctx, bindings, fmt.Sprintf("Allow the agent to download %s into %s", parsed.Redacted(), target)); err != nil {
			return nil, err
		}

		result, err := network.Download(parsed.String(), target, network.DownloadOptions{
			Timeout:        time.Duration(args.TimeoutSeconds) * time.Second,
			MaxBytes:       args.MaxBytes,
			ExpectedSha256: args.ExpectedSha256,
		})
		if err != nil {
			return nil, err
		}

		return &DownloadFileToolOutput{
			URL:      parsed.Redacted(),
			File:     target,
			Bytes:    result.Bytes,
			SHA256:   result.SHA256,
			Verified: result.Verified,
		}, nil
	},
	[]string{tool.NetworkingCategory, tool.FilesystemCategory},
)
