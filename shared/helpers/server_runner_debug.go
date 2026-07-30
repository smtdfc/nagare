//go:build debug

package helpers

import (
	"os"
	"os/exec"
)

func TryStartServer(port, publicKey string) error {
	cmd := exec.Command("go", "run", "../server/main.go", "--port", port, "--pubkey", publicKey)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
