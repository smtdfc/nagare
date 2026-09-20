package declarations

import (
	"time"

	"github.com/smtdfc/nagare/core/context"
	"github.com/smtdfc/nagare/core/tool"
)

type TimeToolInput struct {
	Timezone string `json:"timezone,omitempty"`
}

type TimeToolOutput struct {
	Datetime string `json:"datetime"`
	UnixTime int64  `json:"unix_time"`
}

var TimeTool = tool.DefineTool(
	"time_tool",
	"Get current date and time",
	func(ctx *context.ExecuteContext, args *TimeToolInput, _ tool.Bindings) (*TimeToolOutput, error) {
		loc := time.Local

		if args != nil && args.Timezone != "" {
			if parsedLoc, err := time.LoadLocation(args.Timezone); err == nil {
				loc = parsedLoc
			}
		}

		now := time.Now().In(loc)

		return &TimeToolOutput{
			Datetime: now.Format("2006-01-02 15:04:05 MST"),
			UnixTime: now.Unix(),
		}, nil
	},
)
