//go:build !debug

package helpers

import (
	"os/exec"

	"github.com/smtdfc/nagare/shared/paths"
)

func TryStartServer(port, publicKey string) error {
	cmd := exec.Command(paths.ServerBinFile, "--port", port, "--pubkey", publicKey)
	return cmd.Run()
}
