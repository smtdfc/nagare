//go:build debug
// +build debug

package helpers

import (
	"fmt"
	"os"
	"os/exec"
)

package helpers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func TryStartServer(port string, publicKey string) error {
	fmt.Println("Running server in debug mode")
	cmd := exec.Command("go", "run", "../server/main.go", "--port", port, "--pubkey", publicKey)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("Cannot start server. Cause: %w", err)
	}

	go func() {
		cmd.Wait()
	}()

	return nil
}