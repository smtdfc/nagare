package declarations

import (
	"fmt"

	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/tool"
	"github.com/smtdfc/nagare/cross-platform/media"
)

type SetVolumeArgs struct {
	Volume int `json:"volume" jsonschema:"description=Target system volume from 0 to 100"`
}

var set_volume = tool.DeclareTool(
	"set_volume",
	"Sets the system master volume (0-100).",
	func(ctx domains.AgentContext, args SetVolumeArgs) (any, error) {
		if args.Volume < 0 || args.Volume > 100 {
			return nil, fmt.Errorf("volume must be between 0 and 100")
		}

		if err := media.SetVolume(args.Volume); err != nil {
			return nil, err
		}

		return map[string]any{
			"status": "success",
			"volume": args.Volume,
		}, nil
	},
	domains.STATIC_TOOL,
	domains.NO_GROUP,
)
