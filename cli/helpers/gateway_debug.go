//go:build debug
// +build debug

package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func TryStartGateway(isDebugMode bool, host string) error {
	var cmd *exec.Cmd
	workingDir, _ := os.Getwd()
	gatewayDir := filepath.Join(workingDir, "../gateway")

	if isDebugMode {

		wireCmd := exec.Command("dix", "wire", ".", "--workspace")
		wireCmd.Dir = gatewayDir
		wireCmd.Stdout = os.Stdout
		wireCmd.Stderr = os.Stderr
		wireCmd.Stdin = os.Stdin

		if err := wireCmd.Run(); err != nil {
			return fmt.Errorf("dix wire failed: %w", err)
		}

		cmd = exec.Command("go", "run", "-race", "main.go")
		cmd.Dir = gatewayDir
		cmd.Env = append(os.Environ(), "NAGARE_GATEWAY_MODE=debug")
	}

	if cmd == nil {
		return fmt.Errorf("gateway command is not initialized")
	}

	if host != "" {
		cmd.Env = append(cmd.Environ(), fmt.Sprintf("NAGARE_GATEWAY_HOST=%s", host))
	}

	publicKey, _, err := GetRSAKey()
	if err != nil {
		return err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Env = append(cmd.Environ(), fmt.Sprintf("NAGARE_GATEWAY_PUBLIC_KEY=%s", publicKey))

	return runGateway(cmd)
}
