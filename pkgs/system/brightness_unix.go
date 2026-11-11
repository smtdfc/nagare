//go:build !windows && !darwin

package system

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// brightnessGet reads the current brightness via brightnessctl's
// machine-readable info output.
func brightnessGet() (int, error) {
	out, err := exec.Command("brightnessctl", "info", "--machine-readable").Output()
	if err != nil {
		return 0, fmt.Errorf("brightnessctl is required for brightness control on this system (install it with your package manager): %w", err)
	}

	// Machine-readable format: <name>,<device>,<current>,<max>,<percent>,...
	fields := strings.Split(strings.TrimSpace(string(out)), ",")
	for _, field := range fields {
		trimmed := strings.TrimSpace(field)
		if strings.HasSuffix(trimmed, "%") {
			percent, err := strconv.Atoi(strings.TrimSuffix(trimmed, "%"))
			if err != nil {
				return 0, fmt.Errorf("failed to parse brightnessctl output %q: %w", trimmed, err)
			}
			return percent, nil
		}
	}

	return 0, fmt.Errorf("failed to parse brightnessctl output: %q", string(out))
}

// brightnessSet applies a brightness percentage through brightnessctl.
func brightnessSet(value int) error {
	cmd := exec.Command("brightnessctl", "set", "--quiet", fmt.Sprintf("%d%%", value))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set brightness via brightnessctl: %w", err)
	}
	return nil
}
