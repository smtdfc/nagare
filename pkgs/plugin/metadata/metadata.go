package metadata

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"
)

var validPackageNameRegex = regexp.MustCompile(`^[a-zA-Z0-9]+([a-zA-Z0-9_.-]*[a-zA-Z0-9]+)?$`)

type PluginBinaryMetadata map[string]string
type PluginFeatures []string

type PluginMetadata struct {
	PackageName string               `json:"package_name"`
	Name        string               `json:"name"`
	Version     string               `json:"version"`
	ApiVersion  string               `json:"api_version"`
	Author      string               `json:"author"`
	Bin         PluginBinaryMetadata `json:"bin"`
	Features    PluginFeatures       `json:"features"`
}

func (m *PluginMetadata) Validate() error {
	trimmedPackageName := strings.TrimSpace(m.PackageName)
	if trimmedPackageName == "" {
		return errors.New("plugin metadata 'package_name' is required")
	}
	if strings.Contains(trimmedPackageName, "..") || strings.Contains(trimmedPackageName, "/") || strings.Contains(trimmedPackageName, "\\") {
		return errors.New("plugin metadata 'package_name' contains invalid characters or path traversal")
	}
	if !validPackageNameRegex.MatchString(trimmedPackageName) {
		return errors.New("plugin metadata 'package_name' must contain only alphanumeric characters, dots, underscores, or hyphens")
	}

	if strings.TrimSpace(m.Name) == "" {
		return errors.New("plugin metadata 'name' is required")
	}
	if strings.TrimSpace(m.Version) == "" {
		return errors.New("plugin metadata 'version' is required")
	}
	if strings.TrimSpace(m.ApiVersion) == "" {
		return errors.New("plugin metadata 'api_version' is required")
	}
	if strings.TrimSpace(m.Author) == "" {
		return errors.New("plugin metadata 'author' is required")
	}

	for osName, binPath := range m.Bin {
		trimmedBin := strings.TrimSpace(binPath)
		if trimmedBin == "" {
			return errors.New("plugin metadata 'bin' path for " + osName + " cannot be empty")
		}
		cleanBin := filepath.Clean(trimmedBin)
		normalized := strings.ReplaceAll(cleanBin, "\\", "/")
		if filepath.IsAbs(cleanBin) || strings.HasPrefix(normalized, "/") || strings.HasPrefix(normalized, "../") || normalized == ".." || strings.Contains(normalized, "/../") {
			return errors.New("plugin metadata 'bin' path for " + osName + " is invalid or attempts path traversal")
		}
	}

	return nil
}

func (m *PluginMetadata) GetFeatures() string {
	return strings.Join(m.Features, ",")
}
