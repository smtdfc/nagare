package system

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
