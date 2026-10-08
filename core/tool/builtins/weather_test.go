package declarations

import (
	"testing"

	"github.com/smtdfc/nagare/core/tool"
)

func TestWeatherTool_Metadata(t *testing.T) {
	// Verify WeatherTool metadata
	if WeatherTool.GetName() != "weather_tool" {
		t.Errorf("expected 'weather_tool', got '%s'", WeatherTool.GetName())
	}
	if len(WeatherTool.GetCategories()) != 1 || WeatherTool.GetCategories()[0] != tool.WeatherCategory {
		t.Errorf("expected WeatherCategory, got %v", WeatherTool.GetCategories())
	}
}
