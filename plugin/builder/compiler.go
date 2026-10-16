package builder

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	SourceLangGo   = "go"
	SourceLangRust = "rust"
)

func detectSourceLanguage(dir string) (string, error) {
	_, goErr := os.Stat(filepath.Join(dir, "go.mod"))
	_, cargoErr := os.Stat(filepath.Join(dir, "Cargo.toml"))

	hasGo := goErr == nil
	hasRust := cargoErr == nil

	switch {
	case hasGo && hasRust:
		return "", fmt.Errorf("ambiguous plugin source: both go.mod and Cargo.toml exist (keep only one)")
	case hasRust:
		return SourceLangRust, nil
	case hasGo:
		return SourceLangGo, nil
	default:
		return "", fmt.Errorf("no plugin source found: neither go.mod nor Cargo.toml exists")
	}
}

func compileSource(cwd, binFile, lang string) error {
	switch lang {
	case SourceLangRust:
		return compileRustSource(cwd, binFile)
	case SourceLangGo:
		return compileGoSource(cwd, binFile)
	default:
		return fmt.Errorf("unsupported source language %q", lang)
	}
}

func streamStderr(cmd *exec.Cmd, tag string) error {
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stderrPipe)
	// Cargo prints long JSON-ish lines; allow up to 1 MiB per line.
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("%s[%s]%s %s\n", ColorCyan, tag, ColorReset, line)
	}

	return cmd.Wait()
}

func compileGoSource(cwd, binFile string) error {
	cmd := exec.Command("go", "build", "-o", binFile, ".")
	cmd.Dir = cwd
	return streamStderr(cmd, "GO")
}

func compileRustSource(cwd, binFile string) error {
	binName, err := cargoBinaryName(filepath.Join(cwd, "Cargo.toml"))
	if err != nil {
		return err
	}

	cmd := exec.Command("cargo", "build", "--release")
	cmd.Dir = cwd
	if err := streamStderr(cmd, "RUST"); err != nil {
		return fmt.Errorf("cargo build failed: %w", err)
	}

	built, err := locateCargoBinary(cwd, binName)
	if err != nil {
		return err
	}

	if err := copyBinary(built, binFile); err != nil {
		return fmt.Errorf("failed to stage rust binary: %w", err)
	}
	return nil
}

func locateCargoBinary(cwd, binName string) (string, error) {
	candidates := []string{binName}
	if underscored := strings.ReplaceAll(binName, "-", "_"); underscored != binName {
		candidates = append(candidates, underscored)
	}

	suffix := ""
	if runtime.GOOS == "windows" {
		suffix = ".exe"
	}

	for _, name := range candidates {
		path := filepath.Join(cwd, "target", "release", name+suffix)
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path, nil
		}
	}

	return "", fmt.Errorf("cargo build succeeded but release binary %q was not found in target/release", binName)
}

func copyBinary(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		_ = in.Close()
	}()

	if err := os.MkdirAll(filepath.Dir(dst), 0775); err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	if runtime.GOOS != "windows" {
		if err := os.Chmod(dst, 0755); err != nil {
			return err
		}
	}
	return nil
}

func cargoBinaryName(manifestPath string) (string, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return "", err
	}

	var packageName, binName, section string
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(stripTomlComment(line))
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			inner := strings.Trim(strings.TrimSpace(trimmed), "[]")
			switch inner {
			case "package":
				section = "package"
			case "bin":
				section = "bin"
			default:
				section = ""
			}
			continue
		}

		key, value, ok := splitTomlKeyValue(trimmed)
		if !ok || key != "name" {
			continue
		}
		name := unquoteTomlString(value)
		if name == "" {
			continue
		}
		switch section {
		case "bin":
			if binName == "" {
				binName = name
			}
		case "package":
			if packageName == "" {
				packageName = name
			}
		}
	}

	if binName != "" {
		return binName, nil
	}
	if packageName != "" {
		return packageName, nil
	}
	return "", fmt.Errorf("could not determine binary name from %s (missing [package] name)", manifestPath)
}

func stripTomlComment(line string) string {
	var b strings.Builder
	inSingle, inDouble := false, false
	for _, r := range line {
		switch r {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			}
			b.WriteRune(r)
		case '"':
			if !inSingle {
				inDouble = !inDouble
			}
			b.WriteRune(r)
		case '#':
			if !inSingle && !inDouble {
				return b.String()
			}
			b.WriteRune(r)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func splitTomlKeyValue(line string) (key, value string, ok bool) {
	idx := strings.Index(line, "=")
	if idx < 0 {
		return "", "", false
	}
	key = strings.TrimSpace(line[:idx])
	value = strings.TrimSpace(line[idx+1:])
	if key == "" || value == "" {
		return "", "", false
	}
	return key, value, true
}

func unquoteTomlString(value string) string {
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'') {
			return value[1 : len(value)-1]
		}
	}
	return ""
}
