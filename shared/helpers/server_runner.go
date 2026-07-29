//go:build !debug
// +build !debug

package helpers

import (
	"fmt"
	"os/exec"

	"github.com/smtdfc/nagare/shared/paths"
)

func TryStartServer(port string, publicKey string) error {
	cmd := exec.Command(paths.ServerBinFile, "--port", port, "--pubkey", publicKey)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("Cannot start server. Cause: %w", err)
	}

	return nil
}
