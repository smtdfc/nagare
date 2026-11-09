package system

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const MaxFileToolBytes = 1 << 20

type UserDirectories struct {
	Home      string `json:"home"`
	Desktop   string `json:"desktop"`
	Downloads string `json:"downloads"`
	Documents string `json:"documents"`
	Pictures  string `json:"pictures"`
	Music     string `json:"music"`
	Videos    string `json:"videos"`
	Temp      string `json:"temp"`
	Config    string `json:"config"`
	Cache     string `json:"cache"`
}

func ResolveFilePath(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	return filepath.Clean(path), nil
}

func GetUserDirectories() (UserDirectories, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return UserDirectories{}, fmt.Errorf("failed to resolve user home directory: %w", err)
	}

	directories := UserDirectories{
		Home:      home,
		Downloads: filepath.Join(home, "Downloads"),
		Documents: filepath.Join(home, "Documents"),
		Pictures:  filepath.Join(home, "Pictures"),
		Music:     filepath.Join(home, "Music"),
		Videos:    filepath.Join(home, "Videos"),
		Temp:      os.TempDir(),
		Config:    userConfigDir(home),
		Cache:     userCacheDir(home),
	}

	if runtime.GOOS == "linux" || runtime.GOOS == "freebsd" || runtime.GOOS == "openbsd" || runtime.GOOS == "netbsd" {
		directories.Desktop = xdgDirectory("XDG_DESKTOP_DIR", filepath.Join(home, "Desktop"), home)
		directories.Downloads = xdgDirectory("XDG_DOWNLOAD_DIR", directories.Downloads, home)
		directories.Documents = xdgDirectory("XDG_DOCUMENTS_DIR", directories.Documents, home)
		directories.Pictures = xdgDirectory("XDG_PICTURES_DIR", directories.Pictures, home)
		directories.Music = xdgDirectory("XDG_MUSIC_DIR", directories.Music, home)
		directories.Videos = xdgDirectory("XDG_VIDEOS_DIR", directories.Videos, home)
	} else {
		directories.Desktop = filepath.Join(home, "Desktop")
	}

	return directories, nil
}

func xdgDirectory(variable string, fallback string, home string) string {
	value := os.Getenv(variable)
	if value == "" {
		return fallback
	}
	if value == "$HOME" {
		return home
	}
	if strings.HasPrefix(value, "$HOME") {
		return filepath.Join(home, strings.TrimPrefix(value, "$HOME/"))
	}
	return filepath.Clean(value)
}

func userConfigDir(home string) string {
	if value, err := os.UserConfigDir(); err == nil {
		return value
	}
	return filepath.Join(home, ".config")
}

func userCacheDir(home string) string {
	if value, err := os.UserCacheDir(); err == nil {
		return value
	}
	return filepath.Join(home, ".cache")
}

// SanitizeFilename strips path separators from a filename so the result can
// be safely joined inside its target directory.
func SanitizeFilename(filename string) string {
	name := filepath.Base(strings.TrimSpace(filename))
	if name == "." || name == "/" || name == string(filepath.Separator) {
		return ""
	}
	return name
}

// randomSuffixHex returns a short random hex string for generated filenames.
func randomSuffixHex(n int) string {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(buf)
}

// DeriveDownloadPath resolves the final absolute path for a downloaded file
// based on the URL path, an optional destination (either a directory or a
// file path), and an optional filename:
//   - destination unset: the filename is joined into the caller's working
//     directory of choice by the caller (pass a dir here or set dest)
//   - destination is an existing directory: filename (or URL basename)
//   - destination is an existing file: destination as-is
//   - destination missing with an extension: treated as a file path
//   - destination missing without an extension: treated as a directory
func DeriveDownloadPath(urlPath string, destination string, filename string) (string, error) {
	name := SanitizeFilename(filename)
	if name == "" {
		name = SanitizeFilename(filepath.Base(urlPath))
	}
	if name == "" {
		name = "download-" + randomSuffixHex(4)
	}

	destination, err := ResolveFilePath(destination)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(destination)
	switch {
	case err == nil && info.IsDir():
		return filepath.Join(destination, name), nil
	case err == nil:
		// destination is an existing file path, overwrite it wholesale
		return destination, nil
	default:
		if filepath.Ext(destination) != "" {
			// a non-existing path with an extension is a file target
			if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
				return "", fmt.Errorf("failed to create destination directory: %w", err)
			}
			return destination, nil
		}
		if err := os.MkdirAll(destination, 0o755); err != nil {
			return "", fmt.Errorf("failed to create destination directory: %w", err)
		}
		return filepath.Join(destination, name), nil
	}
}
