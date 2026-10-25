package declarations

import (
	"context"
	"strings"
	"testing"
	"time"

	core_context "github.com/smtdfc/nagare/core/context"
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

func TestWeatherTool_CancelledContext(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	cancel()

	ctx := &core_context.ExecuteContext{Context: parent}
	_, err := WeatherTool.Execute(ctx, `{"lat":21.0,"lng":105.0}`)
	if err == nil {
		t.Fatalf("expected request to fail with a cancelled context")
	}
	if !strings.Contains(err.Error(), "context") {
		t.Errorf("expected a context error, got %q", err.Error())
	}
}

func TestWeatherTool_ExpiredDeadline(t *testing.T) {
	parent, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()

	ctx := &core_context.ExecuteContext{Context: parent}
	_, err := WeatherTool.Execute(ctx, `{"lat":21.0,"lng":105.0}`)
	if err == nil {
		t.Fatalf("expected request to fail with an expired deadline")
	}
}
