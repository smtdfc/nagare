package system

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

// TrashMove moves a file or directory to the platform's trash:
//   - Linux: freedesktop.org trash spec (~/.local/share/Trash or $XDG_DATA_HOME)
//   - macOS: ~/.Trash
//   - Windows: Recycle Bin (via PowerShell's VisualBasic FileSystem API)
//
// The original absolute path and deletion timestamp are preserved so trash
// manager UIs can offer restoration.
func TrashMove(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	path, err := ResolveFilePath(path)
	if err != nil {
		return "", err
	}
	if _, err := os.Lstat(path); err != nil {
		return "", fmt.Errorf("path does not exist: %w", err)
	}

	switch runtime.GOOS {
	case "linux", "freebsd", "netbsd", "openbsd":
		return trashUnix(path)
	case "darwin":
		return trashDarwin(path)
	case "windows":
		return trashWindows(path)
	default:
		return "", fmt.Errorf("trash is not supported on %s", runtime.GOOS)
	}
}

// trashUnix implements the freedesktop.org trash specification for the home
// trash on Unix systems.
func trashUnix(path string) (string, error) {
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to resolve home directory: %w", err)
		}
		dataHome = filepath.Join(home, ".local", "share")
	}

	filesDir := filepath.Join(dataHome, "Trash", "files")
	infoDir := filepath.Join(dataHome, "Trash", "info")
	if err := os.MkdirAll(filesDir, 0o700); err != nil {
		return "", fmt.Errorf("failed to create trash files directory: %w", err)
	}
	if err := os.MkdirAll(infoDir, 0o700); err != nil {
		return "", fmt.Errorf("failed to create trash info directory: %w", err)
	}

	base := filepath.Base(path)
	if _, err := os.Lstat(filepath.Join(filesDir, base)); err == nil {
		base = fmt.Sprintf("%s.%d%s", base, time.Now().UnixNano(), filepath.Ext(base))
	}

	target := filepath.Join(filesDir, base)
	if err := moveTo(target, path); err != nil {
		return "", err
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}

	infoFile := filepath.Join(infoDir, base+".trashinfo")
	infoContent := fmt.Sprintf("[Trash Info]\nPath=%s\nDeletionDate=%s\n",
		escapeTrashInfo(abs), time.Now().Format("2006-01-02T15:04:05"))
	if err := os.WriteFile(infoFile, []byte(infoContent), 0o600); err != nil {
		// The file is already swapped into the trash; do not fail the move.
		return filesDir, fmt.Errorf("failed to write trash metadata for %s: %w", path, err)
	}

	return target, nil
}

// trashDarwin moves paths into ~/.Trash, which trash managers index natively.
func trashDarwin(path string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve home directory: %w", err)
	}

	target := filepath.Join(home, ".Trash", filepath.Base(path))
	if _, err := os.Lstat(target); err == nil {
		target = fmt.Sprintf("%s.%d%s", target, time.Now().UnixNano(), filepath.Ext(path))
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
		return "", fmt.Errorf("failed to create trash directory: %w", err)
	}
	if err := moveTo(target, path); err != nil {
		return "", err
	}
	return target, nil
}

// trashWindows deletes via the Recycle Bin using the portable PowerShell
// bridge: Microsoft.VisualBasic.FileIO.FileSystem with SendToRecycleBin.
func trashWindows(path string) (string, error) {
	verb := "DeleteFile"
	if info, err := os.Lstat(path); err != nil {
		return "", fmt.Errorf("path does not exist: %w", err)
	} else if info.IsDir() {
		verb = "DeleteDirectory"
	}

	psPath := escapePowerShell(path)
	script := fmt.Sprintf(
		"Add-Type -AssemblyName Microsoft.VisualBasic; [Microsoft.VisualBasic.FileIO.FileSystem]::%s(%s, 'OnlyErrorDialogs', 'SendToRecycleBin')",
		verb, psPath,
	)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to move '%s' to the recycle bin: %w", path, err)
	}

	return "Recycle Bin", nil
}

// moveTo renames src to target, falling back to a copy+delete for files that
// cross filesystem boundaries.
func moveTo(target string, src string) error {
	if err := os.Rename(src, target); err == nil {
		return nil
	} else if !isCrossDevice(err) {
		return fmt.Errorf("failed to move '%s' to trash: %w", src, err)
	}

	info, err := os.Lstat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("failed to move directory '%s' to trash across filesystems; copy it to the same device first", src)
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_EXCL, info.Mode().Perm())
	if err != nil {
		return fmt.Errorf("failed to create trash copy of '%s': %w", src, err)
	}

	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		os.Remove(target)
		return fmt.Errorf("failed to copy '%s' to trash: %w", src, err)
	}
	if err := out.Close(); err != nil {
		os.Remove(target)
		return fmt.Errorf("failed to finalize trash copy of '%s': %w", src, err)
	}

	if err := os.Remove(src); err != nil {
		return fmt.Errorf("failed to remove original '%s' after moving to trash: %w", src, err)
	}

	return nil
}

func escapePowerShell(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}

// isCrossDevice reports whether the rename failure was caused by a
// cross-filesystem move (EXDEV).
func isCrossDevice(err error) bool {
	return errors.Is(err, syscall.EXDEV)
}

func escapeTrashInfo(path string) string {
	replacer := map[rune]string{'&': "&amp;", '<': "&lt;", '>': "&gt;"}
	out := make([]rune, 0, len(path))
	for _, r := range path {
		if s, ok := replacer[r]; ok {
			out = append(out, []rune(s)...)
			continue
		}
		out = append(out, r)
	}
	return string(out)
}
