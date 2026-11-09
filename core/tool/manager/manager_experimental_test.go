package manager

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/smtdfc/nagare/core/logger"
)

func testLoggerWithBuffer() (*logger.BaseLogger, *bytes.Buffer) {
	buffer := &bytes.Buffer{}
	created := &logger.BaseLogger{
		Logger: *slog.New(slog.NewTextHandler(buffer, nil)),
	}
	return created, buffer
}

func TestWarnExperimentalComputerUse(t *testing.T) {
	testLogger, buffer := testLoggerWithBuffer()

	warnExperimentalComputerUse(testLogger, "mouse_action_tool")
	warnExperimentalComputerUse(testLogger, "mouse_action_tool")

	output := buffer.String()
	if !strings.Contains(output, "Experimental tools registered") {
		t.Errorf("expected a one-time warning listing experimental tools, got: %s", output)
	}
	if !strings.Contains(output, "computer_use") {
		t.Errorf("expected the warning to identify the computer_use category, got: %s", output)
	}
	if strings.Count(output, "Experimental tools registered") != 1 {
		t.Errorf("expected the registration warning to log exactly once, got output:\n%s", output)
	}
	if strings.Count(output, "Experimental feature used") != 2 {
		t.Errorf("expected a usage warning per call, got output:\n%s", output)
	}
}

func TestHasComputerUseCategory(t *testing.T) {
	if !hasComputerUseCategory([]string{"system", "computer_use"}) {
		t.Errorf("expected computer_use to be detected")
	}
	if hasComputerUseCategory([]string{"system", "filesystem"}) {
		t.Errorf("expected non-computer-use categories to be ignored")
	}
	if hasComputerUseCategory(nil) {
		t.Errorf("expected nil categories to be ignored")
	}
}
