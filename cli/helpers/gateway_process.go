package helpers

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/smtdfc/nagare/pkgs/paths"
)

func runGateway(cmd *exec.Cmd) error {
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := saveGatewayPID(cmd.Process.Pid); err != nil {
		_ = requestGatewayShutdown(cmd.Process)
		_ = cmd.Wait()
		return fmt.Errorf("failed to save gateway PID: %w", err)
	}

	fmt.Printf("Starting Nagare gateway with PID %d...\n", cmd.Process.Pid)
	err := cmd.Wait()
	if removeErr := removeGatewayPID(); err == nil {
		err = removeErr
	}

	return err
}

func saveGatewayPID(pid int) error {
	return os.WriteFile(paths.GatewayPIDFile, []byte(strconv.Itoa(pid)), 0644)
}

func removeGatewayPID() error {
	err := os.Remove(paths.GatewayPIDFile)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return err
}

func TryStopGateway() (bool, error) {
	data, err := os.ReadFile(paths.GatewayPIDFile)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to read gateway PID file: %w", err)
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil || pid <= 0 {
		_ = removeGatewayPID()
		return false, fmt.Errorf("invalid gateway PID in %s", paths.GatewayPIDFile)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		_ = removeGatewayPID()
		return false, nil
	}

	if err := requestGatewayShutdown(process); err != nil {
		return false, fmt.Errorf("failed to stop gateway process %d: %w", pid, err)
	}

	return true, nil
}
