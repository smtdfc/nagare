//go:build debug

package helpers

import (
	"os"
	"os/exec"
	"strings"
)

func TryStartServer(port, publicKey string) error {
	cleanKey := strings.ReplaceAll(publicKey, "\n", " ")
	script := "cd ../server && dix wire && go run . --port " + port + " --pubkey \"" + cleanKey + "\""

	cmd := exec.Command("sh", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
