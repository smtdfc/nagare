//go:build !debug
// +build !debug

package helpers

import (
	"fmt"
	"os"
	"os/exec"
)

func TryStartGateway(isDebugMode bool) error {
	var cmd *exec.Cmd

	cmd = exec.Command("nagare-gateway")
	cmd.Env = append(os.Environ(), "NAGARE_GATEWAY_MODE=prod")

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
