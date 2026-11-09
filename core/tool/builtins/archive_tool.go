package declarations

import (
	"fmt"
	"strings"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/system"
)

type ArchiveToolInput struct {
	Action      string   `json:"action" jsonschema_description:"Operation to perform: 'create' packs sources into archive_path, 'extract' unpacks archive_path into destination."`
	ArchivePath string   `json:"archive_path" jsonschema_description:"Path of the archive. For create, the extension defines the format (.zip, .tar, .tar.gz, .tgz); the file must not exist yet. For extract, the archive to unpack."`
	Sources     []string `json:"sources,omitempty" jsonschema_description:"Files or directories to pack when action is 'create'. Directories are included recursively, preserving their structure."`
	Destination string   `json:"destination,omitempty" jsonschema_description:"Target directory to unpack into when action is 'extract'. Defaults to the directory containing the archive."`
	Confirm     bool     `json:"confirm" jsonschema_description:"Must be true to allow the archive operation."`
}

type ArchiveToolOutput struct {
	Action      string   `json:"action" jsonschema_description:"The operation that was performed."`
	ArchivePath string   `json:"archive_path"`
	Destination string   `json:"destination" jsonschema_description:"The target directory for extract, empty for create."`
	Format      string   `json:"format" jsonschema_description:"Detected archive format: zip, tar, or tar.gz."`
	Sources     []string `json:"sources" jsonschema_description:"The packed source paths for create, empty for extract."`
	Files       int      `json:"files" jsonschema_description:"Number of source files packed for create, empty for extract."`
}

var ArchiveTool = tool.DefineTool(
	"archive_tool",
	"Create or extract zip and tar archives: pack files and directories into .zip/.tar/.tar.gz, or safely unpack an archive into a directory with zip-slip and archive bomb protection.",
	func(ctx *context.ExecuteContext, args *ArchiveToolInput, bindings tool.Bindings) (*ArchiveToolOutput, error) {
		if args == nil || !args.Confirm {
			return nil, fmt.Errorf("archive operation requires confirm=true")
		}

		action := strings.ToLower(strings.TrimSpace(args.Action))
		if action != "create" && action != "extract" {
			return nil, fmt.Errorf("action is required ('create' or 'extract')")
		}
		if args.ArchivePath == "" {
			return nil, fmt.Errorf("archive_path is required")
		}

		format := system.DetectArchiveFormat(args.ArchivePath)
		if format == system.ArchiveFormatUnknown {
			return nil, fmt.Errorf("unsupported archive format: '%s' (supported: .zip, .tar, .tar.gz, .tgz)", extOf(args.ArchivePath))
		}

		if err := requireApproval(ctx, bindings, fmt.Sprintf("Allow the agent to %s archive %s", action, args.ArchivePath)); err != nil {
			return nil, err
		}

		switch action {
		case "create":
			if len(args.Sources) == 0 {
				return nil, fmt.Errorf("at least one source is required when creating an archive")
			}
			if err := system.ArchiveCreate(args.Sources, args.ArchivePath); err != nil {
				return nil, err
			}
			return &ArchiveToolOutput{
				Action:      "create",
				ArchivePath: args.ArchivePath,
				Format:      string(format),
				Sources:     args.Sources,
				Files:       countPackedFiles(args.Sources),
			}, nil

		default:
			destination, err := system.ArchiveExtract(args.ArchivePath, args.Destination)
			if err != nil {
				return nil, err
			}
			return &ArchiveToolOutput{
				Action:      "extract",
				ArchivePath: args.ArchivePath,
				Destination: destination,
				Format:      string(format),
			}, nil
		}
	},
	[]string{tool.FilesystemCategory},
)

func extOf(path string) string {
	lower := strings.ToLower(path)
	if strings.HasSuffix(lower, ".tar.gz") {
		return ".tar.gz"
	}
	return strings.ToLower(path) // lower ext report
}

// countPackedFiles counts the regular files available under the sources so
// the operator can sanity-check the result size.
func countPackedFiles(sources []string) int {
	total := 0
	for _, source := range sources {
		files, err := system.ArchiveSourceFiles(source)
		if err != nil {
			continue
		}
		total += len(files)
	}
	return total
}
