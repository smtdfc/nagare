package declarations

import (
	"context"
	"time"

	"github.com/smtdfc/nagare/core/tool"
)

type LocalTimeArgs struct{}

var get_local_time = tool.DeclareTool(
	"get_local_time",
	"Get the current local date and time of the system.",
	func(ctx context.Context, args LocalTimeArgs) (any, error) {
		now := time.Now()

		return map[string]any{
			"timestamp":    now.Unix(),
			"time":         now.Format("2006-01-02 15:04:05"),
			"time_rfc3339": now.Format(time.RFC3339),
			"timezone":     now.Location().String(),
			"offset":       now.Format("-07:00"),
		}, nil
	},
	tool.STATIC_TOOL,
	tool.NO_GROUP,
)
