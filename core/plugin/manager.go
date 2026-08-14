package plugin

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/persistence/database/mappers"
	"github.com/smtdfc/nagare/core/persistence/database/models"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
	"github.com/smtdfc/nagare/shared/paths"
	nagare_plugin "github.com/smtdfc/nagare/shared/plugin"
)

type PluginManager struct {
	pluginRepo *repositories.PluginRepository
}

func (p *PluginManager) ReadMetadata(pluginMetadataPath string) (*nagare_plugin.PluginMetadata, error) {
	metadata := &nagare_plugin.PluginMetadata{}
	raw, err := os.ReadFile(pluginMetadataPath)
	if err != nil {
		PluginLogger.Error("failed to read plugin metadata", "error", err)
		return nil, custom_errors.NewPluginError("failed to read plugin metadata")
	}
	if err := json.Unmarshal(raw, metadata); err != nil {
		PluginLogger.Error("failed to unmarshal plugin metadata", "error", err)
		return nil, custom_errors.NewPluginError("failed to unmarshal plugin metadata")
	}
	return metadata, nil
}

func (p *PluginManager) GetPluginDirName(metadata *nagare_plugin.PluginMetadata) string {
	return fmt.Sprintf("%s_%s_%s", metadata.ID, metadata.Version, runtime.GOARCH)
}

func (p *PluginManager) copyPlugin(pluginPath, pluginDirPath string) error {
	if err := os.MkdirAll(pluginDirPath, 0755); err != nil {
		return err
	}

	srcFile, err := os.Open(pluginPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFileName := filepath.Base(pluginPath)
	destFilePath := filepath.Join(pluginDirPath, destFileName)

	dstFile, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return dstFile.Sync()
}

func (p *PluginManager) RegisterPlugin(pluginPath string) error {
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return custom_errors.NewPluginError("plugin path not found")
	}

	pluginMetatadaFile := path.Join(pluginPath, "metadata.json")
	if _, err := os.Stat(pluginMetatadaFile); os.IsNotExist(err) {
		return custom_errors.NewPluginError("plugin metadata file not found")
	}

	metadata, err := p.ReadMetadata(pluginMetatadaFile)
	if err != nil {
		return err
	}

	currentArch := runtime.GOARCH
	if !metadata.SupportsArchitecture(currentArch) {
		PluginLogger.Error("plugin architecture mismatch", "current_arch", currentArch)
		return custom_errors.NewPluginError("plugin architecture mismatch")
	}

	pluginDirName := p.GetPluginDirName(metadata)
	pluginDirPath := path.Join(paths.PluginDir, pluginDirName)
	if err := p.copyPlugin(pluginPath, pluginDirPath); err != nil {
		return err
	}
	binFile, err := metadata.GetBinForArchitecture(runtime.GOARCH)
	if err != nil {
		return err
	}

	plugin := &models.Plugin{
		PluginID:   metadata.ID,
		Name:       metadata.Name,
		Version:    metadata.Version,
		ApiVersion: metadata.ApiVersion,
		Author:     metadata.Author,
		Bin:        path.Join(pluginDirPath, binFile.Path),
	}

	err = p.pluginRepo.CreatePlugin(plugin)
	if err != nil {
		return custom_errors.NewPluginError("failed to create plugin")
	}

	PluginLogger.Info("plugin installed", "plugin_dir", pluginDirPath)
	return nil
}

func (p *PluginManager) GetAllPlugins() ([]*domains.PluginInfo, error) {
	mapper := mappers.PluginMapper{}
	plugins, err := p.pluginRepo.GetAllPlugins()
	if err != nil {
		return nil, err
	}

	pluginInfos := make([]*domains.PluginInfo, 0, len(plugins))
	for _, plugin := range plugins {
		pluginInfos = append(pluginInfos, mapper.ToDomain(&plugin))
	}
	return pluginInfos, nil
}

func NewPluginManager(pluginRepo *repositories.PluginRepository) *PluginManager {
	return &PluginManager{
		pluginRepo: pluginRepo,
	}
}
