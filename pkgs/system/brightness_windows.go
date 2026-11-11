//go:build windows

package system

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// runPowerShell executes a single PowerShell expression and captures stdout.
func runPowerShell(script string) (string, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	return string(out), err
}

func brightnessGet() (int, error) {
	out, err := runPowerShell("(Get-CimInstance -Namespace root/WMI -ClassName WmiMonitorBrightness).CurrentBrightness")
	if err != nil {
		return 0, fmt.Errorf("failed to read brightness (desktop WMI monitors only): %w", err)
	}
	value, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return 0, fmt.Errorf("failed to parse brightness output %q: %w", out, err)
	}
	return value, nil
}

func brightnessSet(value int) error {
	_, err := runPowerShell(fmt.Sprintf(
		"Get-CimInstance -Namespace root/WMI -ClassName WmiMonitorBrightnessMethods | Invoke-CimMethod -MethodName WmiSetBrightness -Arguments @{Timeout=1; Brightness=%d}",
		value,
	))
	if err != nil {
		return fmt.Errorf("failed to set brightness (desktop WMI monitors only): %w", err)
	}
	return nil
}
