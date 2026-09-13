package builder

import (
	"bufio"
	"fmt"
	"os/exec"
)

func compileSource(binFile string) error {
	cmd := exec.Command("go", "build", "-a", "-v", "-o", binFile, ".")

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("Failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stderrPipe)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("%s[GO]%s %s\n", ColorCyan, ColorReset, line)
	}

	return cmd.Wait()
}
