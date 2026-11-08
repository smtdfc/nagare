package declarations

import (
	"path/filepath"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/pkgs/paths"
	"github.com/smtdfc/nagare/pkgs/system"
)

type ScreenshotInput struct {
	Index int `json:"index"`
}

type ScreenshotOutput struct {
	File string `json:"file"`
}

var ScreenshotTool = tool.DefineTool(
	"screenshot_tool",
	"Capture a screenshot of the display at the given index and save it as a PNG artifact.",
	func(ctx *context.ExecuteContext, args *ScreenshotInput, bindings tool.Bindings) (*ScreenshotOutput, error) {
		id := uuid.New().String()
		fileName := filepath.Join(paths.ArtifactsDir, id+".png")
		err := system.Screenshot(args.Index, fileName)
		if err != nil {
			return nil, err
		}

		return &ScreenshotOutput{File: fileName}, nil
	},
	[]string{tool.SystemCategory},
)
