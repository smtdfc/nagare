package builder

import (
	"testing"
)

func TestLoggerOutput(t *testing.T) {
	// Verify logger helpers do not panic
	printStep("Testing step")
	printSuccess()
	printError("Testing error %s", "details")
}
