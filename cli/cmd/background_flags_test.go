package cmd

import (
	"testing"

	"github.com/smtdfc/nagare/cli/helpers"
)

func TestBackgroundFlagsRegistered(t *testing.T) {
	startFlag := gatewayStartCmd.Flags().Lookup("background")
	if startFlag == nil {
		t.Fatalf("expected the gateway start command to register --background")
	}
	if startFlag.Shorthand != "b" {
		t.Errorf("expected the --background shorthand 'b', got %q", startFlag.Shorthand)
	}

	webFlag := webCmd.Flags().Lookup("background")
	if webFlag == nil {
		t.Fatalf("expected the web command to register --background")
	}
	if webFlag.Shorthand != "b" {
		t.Errorf("expected the --background shorthand 'b', got %q", webFlag.Shorthand)
	}
}

func TestIsWebDaemonDefaultsToFalse(t *testing.T) {
	t.Setenv("NAGARE_WEB_DAEMON", "")
	if helpers.IsWebDaemon() {
		t.Errorf("expected IsWebDaemon to be false without the daemon env var")
	}
}
