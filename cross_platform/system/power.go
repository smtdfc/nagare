package system

import (
	"fmt"
	"os/exec"
	"runtime"
)

func ExecutePowerCommand(action string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		switch action {
		case "shutdown":
			cmd = exec.Command("shutdown", "/s", "/t", "0")
		case "reboot":
			cmd = exec.Command("shutdown", "/r", "/t", "0")
		case "sleep":
			cmd = exec.Command("powershell", "-Command", "Add-Type -AssemblyName System.Windows.Forms; [System.Windows.Forms.Application]::SetSuspendState('Suspend', $false, $false)")
		}
	case "linux", "darwin":
		switch action {
		case "shutdown":
			cmd = exec.Command("shutdown", "now")
		case "reboot":
			cmd = exec.Command("reboot")
		case "sleep":
			cmd = exec.Command("systemctl", "suspend")
		}
	}

	if cmd == nil {
		return fmt.Errorf("unsupported action or OS")
	}

	return cmd.Run()
}
