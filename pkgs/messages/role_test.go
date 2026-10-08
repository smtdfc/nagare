package messages

import (
	"testing"
)

func TestRoleConstants(t *testing.T) {
	if AGENT != "AGENT" || USER != "USER" || SYSTEM != "SYSTEM" || DEVELOPER != "DEVELOPER" {
		t.Fatalf("unexpected role constant values")
	}
}
