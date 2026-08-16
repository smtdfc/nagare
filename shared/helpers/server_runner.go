package helpers

import (
	"os"
	"os/exec"

	"github.com/smtdfc/nagare/shared/paths"
)

func TryStartServer(publicKey string, debugMode bool) error {
	isRemote := IsRunWithRemoteServer()
	if isRemote {
		return nil
	}

	port, err := GetServerPort()
	if err != nil {
		return err
	}

	if debugMode {
		// cleanKey := strings.ReplaceAll(publicKey, "\n", " ")
		script := "cd ../server && dix wire && go run . --port " + port
		cmd := exec.Command("sh", "-c", script)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = append(os.Environ(), "NAGARE_PUBLIC_KEY="+publicKey)
		return cmd.Run()
	}

	cmd := exec.Command(paths.ServerBinFile, "--port", port, "--pubkey", publicKey)
	return cmd.Run()
}
