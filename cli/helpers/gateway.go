package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func TryStartGateway(isDebugMode bool) error {
	var cmd *exec.Cmd
	if isDebugMode {
		workingDir, _ := os.Getwd()
		cmd = exec.Command("dix", "run", ".", "--workspace")
		cmd.Dir = filepath.Join(workingDir, "../gateway")
		cmd.Env = append(os.Environ(), "NAGARE_GATEWAY_MODE=debug")
	} else {
		cmd = exec.Command("nagare-gateway")
		cmd.Env = append(os.Environ(), "NAGARE_GATEWAY_MODE=prod")
	}

	publicKey, _, err := GetRSAKey()
	if err != nil {
		return err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Env = append(cmd.Environ(), fmt.Sprintf("NAGARE_GATEWAY_PUBLIC_KEY=%s", publicKey))
	return cmd.Run()
}
