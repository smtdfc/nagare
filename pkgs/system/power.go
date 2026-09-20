package system

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type PowerControl struct {
}

func NewPowerControl() *PowerControl {
	return &PowerControl{}
}

func (p *PowerControl) Shutdown() error {
	return p.execute("shutdown")
}

func (p *PowerControl) Restart() error {
	return p.execute("restart")
}

func (p *PowerControl) Suspend() error {
	return p.execute("suspend")
}

func (p *PowerControl) Hibernate() error {
	return p.execute("hibernate")
}

func (p *PowerControl) Logout() error {
	return p.execute("logout")
}

func (p *PowerControl) execute(action string) error {
	action = strings.ToLower(strings.TrimSpace(action))
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		switch action {
		case "shutdown":
			cmd = exec.Command("systemctl", "poweroff")
		case "restart", "reboot":
			cmd = exec.Command("systemctl", "reboot")
		case "suspend", "sleep":
			cmd = exec.Command("systemctl", "suspend")
		case "hibernate":
			cmd = exec.Command("systemctl", "hibernate")
		case "logout":
			cmd = exec.Command("loginctl", "terminate-session", "self")
		}

	case "windows":
		switch action {
		case "shutdown":
			cmd = exec.Command("shutdown", "/s", "/t", "0")
		case "restart", "reboot":
			cmd = exec.Command("shutdown", "/r", "/t", "0")
		case "suspend", "sleep":
			cmd = exec.Command("rundll32.exe", "powrprof.dll,SetSuspendState", "0,1,0")
		case "hibernate":
			cmd = exec.Command("shutdown", "/h")
		case "logout":
			cmd = exec.Command("shutdown", "/l")
		}

	case "darwin": // macOS
		switch action {
		case "shutdown":
			cmd = exec.Command("osascript", "-e", "tell application \"System Events\" to shut down")
		case "restart", "reboot":
			cmd = exec.Command("osascript", "-e", "tell application \"System Events\" to restart")
		case "suspend", "sleep":
			cmd = exec.Command("pmset", "sleepnow")
		case "logout":
			cmd = exec.Command("osascript", "-e", "tell application \"System Events\" to log out")
		case "hibernate":
			return fmt.Errorf("hibernate is not directly supported on macOS via standard CLI")
		}

	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if cmd == nil {
		return fmt.Errorf("invalid action '%s' for OS '%s'", action, runtime.GOOS)
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute %s on %s: %w", action, runtime.GOOS, err)
	}

	return nil
}
