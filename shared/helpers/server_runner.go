package helpers

import (
	"os"
	"os/exec"
	"strings"

	"github.com/smtdfc/nagare/shared/paths"
)

func TryStartServer(port, publicKey string, debugMode bool) error {

	if debugMode {
		cleanKey := strings.ReplaceAll(publicKey, "\n", " ")
		script := "cd ../server && dix wire && go run . --port " + port + " --pubkey \"" + cleanKey + "\""

		cmd := exec.Command("sh", "-c", script)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		return cmd.Start()
	}

	cmd := exec.Command(paths.ServerBinFile, "--port", port, "--pubkey", publicKey)
	return cmd.Run()
}
