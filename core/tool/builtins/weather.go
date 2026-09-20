package declarations

import (
	"context"

	"github.com/smtdfc/nagare/core/tool"
)

type WeatherToolInput struct {
	Lat float32
	Lng float32
}

type WeatherToolOutput struct {
	Temp float32
}

var WeatherTool = tool.DefineTool[WeatherToolInput, WeatherToolOutput](
	"weather_tool",
	"Get weather",
	func(ctx context.Context, args *WeatherToolInput) (*WeatherToolOutput, error) {
		return &WeatherToolOutput{Temp: 34}, nil
	},
)
