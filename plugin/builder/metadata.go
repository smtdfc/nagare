package builder

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/smtdfc/nagare/plugin/metadata"
	"github.com/smtdfc/nagare/shared/helpers"
)

func loadMetadata(metadataFile string) (*metadata.PluginMetadata, error) {
	_, err := os.Stat(metadataFile)
	if errors.Is(err, os.ErrNotExist) {
		return nil, errors.New("Metadata file not exist")
	}

	metadataContent, err := os.ReadFile(metadataFile)
	if err != nil {
		return nil, err
	}

	pluginMetadata, err := helpers.UnmarshalJson[metadata.PluginMetadata](string(metadataContent))
	if err != nil {
		return nil, err
	}

	return pluginMetadata, nil
}

func validateMetadata(pluginMetadata *metadata.PluginMetadata) error {
	return pluginMetadata.Validate()
}

func savePackageMetadata(pkgDir, binFile string, pluginMetadata metadata.PluginMetadata) error {
	relPath, err := filepath.Rel(pkgDir, binFile)
	if err != nil {
		return err
	}

	relPath = filepath.ToSlash(relPath)
	pluginMetadata.Bin = metadata.PluginBinaryMetadata{
		runtime.GOOS: relPath,
	}

	raw, err := helpers.MarshalJson(pluginMetadata)
	if err != nil {
		return err
	}

	newPluginMetadataFile := filepath.Join(pkgDir, "metadata.json")
	return os.WriteFile(newPluginMetadataFile, []byte(raw), 0664)
}
