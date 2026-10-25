package declarations

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/smtdfc/nagare/core/tool"
)

func TestTimeTool_Metadata(t *testing.T) {
	// Verify TimeTool metadata
	if TimeTool.GetName() != "time_tool" {
		t.Errorf("expected name 'time_tool', got '%s'", TimeTool.GetName())
	}
	if len(TimeTool.GetCategories()) != 1 || TimeTool.GetCategories()[0] != tool.TimingCategory {
		t.Errorf("expected TimingCategory, got %v", TimeTool.GetCategories())
	}
	if TimeTool.GetArgsSchema() == "" {
		t.Errorf("expected non-empty schema")
	}
}

func decodeTime(t *testing.T, res string) TimeToolOutput {
	t.Helper()

	var out TimeToolOutput
	if err := json.Unmarshal([]byte(res), &out); err != nil {
		t.Fatalf("failed to decode output %q: %v", res, err)
	}
	return out
}

func TestTimeTool_Execute(t *testing.T) {
	before := time.Now().Add(-time.Minute)
	res, err := TimeTool.Execute(nil, `{}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := decodeTime(t, res)
	if out.Datetime == "" {
		t.Errorf("expected a non-empty datetime")
	}
	got := time.Unix(out.UnixTime, 0)
	if got.Before(before) || got.After(time.Now().Add(time.Minute)) {
		t.Errorf("expected unix time near now, got %s", got)
	}
}

func TestTimeTool_UTCTimezone(t *testing.T) {
	res, err := TimeTool.Execute(nil, `{"timezone":"UTC"}`)
	if err != nil {
		t.Fatalf("unexpected error with UTC: %v", err)
	}

	out := decodeTime(t, res)
	if !strings.HasSuffix(out.Datetime, "UTC") {
		t.Errorf("expected datetime to end with 'UTC', got %q", out.Datetime)
	}
}

func TestTimeTool_InvalidTimezoneFallsBack(t *testing.T) {
	res, err := TimeTool.Execute(nil, `{"timezone":"Not/AZone"}`)
	if err != nil {
		t.Fatalf("expected invalid timezone to fall back without error, got %v", err)
	}

	out := decodeTime(t, res)
	if out.UnixTime == 0 {
		t.Errorf("expected a valid unix time after fallback")
	}
}
