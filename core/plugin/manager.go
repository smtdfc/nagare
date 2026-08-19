package plugin

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/persistence/database/mappers"
	"github.com/smtdfc/nagare/core/persistence/database/models"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
	"github.com/smtdfc/nagare/shared/paths"
	nagare_plugin "github.com/smtdfc/nagare/shared/plugin"
)

const CONNECT_CODE_TTL = 30 * time.Minute

func Unzip(src string, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		filePath := filepath.Join(dest, file.Name)
		if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(filepath.Separator)) {
			continue
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := file.Open()
		if err != nil {
			dstFile.Close()
			return err
		}

		_, err = io.Copy(dstFile, fileInArchive)

		fileInArchive.Close()
		dstFile.Close()

		if err := os.Chmod(filePath, file.Mode()); err != nil {
			return err
		}

		if err != nil {
			return err
		}
	}
	return nil
}

type PluginManager struct {
	connectCodeMgr *ConnectCodeManager
	pluginRepo     *repositories.PluginRepository
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
		PluginLogger.Error("failed to open plugin file", "error", err)
		return custom_errors.NewPluginError("failed to open plugin file")
	}
	defer srcFile.Close()

	destFileName := filepath.Base(pluginPath)
	destFilePath := filepath.Join(pluginDirPath, destFileName)

	dstFile, err := os.Create(destFilePath)
	if err != nil {
		PluginLogger.Error("failed to create destination file", "error", err)
		return custom_errors.NewPluginError("failed to create destination file")
	}
	defer dstFile.Close()
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		PluginLogger.Error("failed to copy plugin file", "error", err)
		return custom_errors.NewPluginError("failed to copy plugin file")
	}

	return dstFile.Sync()
}

func (p *PluginManager) RegisterPlugin(pluginPath string) error {
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return custom_errors.NewPluginError("plugin path not found")
	}

	pluginUUID := uuid.New().String()
	pluginTempDir := filepath.Join(paths.TempDir, pluginUUID)
	err := Unzip(pluginPath, pluginTempDir)
	if err != nil {
		PluginLogger.Error("failed to extract plugin", "error", err)
		return custom_errors.NewPluginError("failed to extract plugin. Please check the file format and try again")
	}
	defer os.RemoveAll(pluginTempDir)

	pluginMetadataFile := filepath.Join(pluginTempDir, "metadata.json")
	if _, err := os.Stat(pluginMetadataFile); os.IsNotExist(err) {
		return custom_errors.NewPluginError("plugin metadata file not found")
	}

	metadata, err := p.ReadMetadata(pluginMetadataFile)
	if err != nil {
		return err
	}

	currentArch := runtime.GOARCH
	if !metadata.SupportsArchitecture(currentArch) {
		PluginLogger.Error("plugin architecture mismatch", "current_arch", currentArch)
		return custom_errors.NewPluginError("plugin architecture mismatch")
	}

	pluginDirName := p.GetPluginDirName(metadata)
	pluginDirPath := filepath.Join(paths.PluginDir, pluginDirName)
	if err := os.MkdirAll(pluginDirPath, 0755); err != nil {
		return err
	}

	if err := Unzip(pluginPath, pluginDirPath); err != nil {
		PluginLogger.Error("failed to extract plugin to destination", "error", err)
		return custom_errors.NewPluginError("failed to extract plugin to destination")
	}

	binFile, err := metadata.GetBinForArchitecture(runtime.GOARCH)
	if err != nil {
		return err
	}
	mapper := &mappers.PluginMapper{}
	plugin := &models.Plugin{
		PluginID:   metadata.ID,
		Name:       metadata.Name,
		Version:    metadata.Version,
		ApiVersion: metadata.ApiVersion,
		Author:     metadata.Author,
		Bin:        filepath.Join(pluginDirPath, binFile.Path),
	}

	err = p.pluginRepo.CreatePlugin(plugin)
	if err != nil {
		return custom_errors.NewPluginError("failed to create plugin")
	}

	PluginLogger.Info("plugin installed", "plugin_dir", pluginDirPath)

	err = p.SpawnPluginProcess(mapper.ToDomain(plugin))
	if err != nil {
		return err
	}

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

func (p *PluginManager) SpawnPluginProcess(plugin *domains.PluginInfo) error {

	binDir := filepath.Dir(plugin.Bin)
	pidFile := filepath.Join(binDir, ".pid")
	connectCode := p.connectCodeMgr.GenerateConnectCode(plugin, CONNECT_CODE_TTL)
	cmd := exec.Command(plugin.Bin)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "NAGARE_PLUGIN_CONNECT_CODE="+connectCode)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		PluginLogger.Error("failed to spawn plugin process", "error", err, "plugin", plugin.PluginID)
		return custom_errors.NewPluginError("failed to spawn plugin process")
	}

	if err := os.WriteFile(pidFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0644); err != nil {
		PluginLogger.Error("failed to write pid file", "error", err, "plugin", plugin.PluginID)
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			PluginLogger.Error("plugin process exited", "error", err, "plugin", plugin.PluginID)
		}
	}()

	return nil
}

func (p *PluginManager) Start() error {
	mapper := &mappers.PluginMapper{}
	plugins, err := p.pluginRepo.GetAllActivePlugins()
	if err != nil {
		PluginLogger.Error("failed to get active plugins", "error", err)
		return custom_errors.NewPluginError("failed to start plugins")
	}

	for _, plugin := range plugins {
		PluginLogger.Info("Starting plugin", "plugin", plugin.PluginID)
		if err := p.SpawnPluginProcess(mapper.ToDomain(&plugin)); err != nil {
			PluginLogger.Error("failed to spawn plugin process", "error", err, "plugin", plugin.PluginID)
			return custom_errors.NewPluginError("failed to start plugins")
		}
	}
	return nil
}

func (p *PluginManager) Stop() error {
	mapper := &mappers.PluginMapper{}
	plugins, err := p.pluginRepo.GetAllActivePlugins()
	if err != nil {
		PluginLogger.Error("failed to get active plugins", "error", err)
		return custom_errors.NewPluginError("failed to stop plugins")
	}

	for _, plugin := range plugins {
		PluginLogger.Info("Stopping plugin", "plugin", plugin.PluginID)
		if err := p.StopPluginProcess(mapper.ToDomain(&plugin)); err != nil {
			PluginLogger.Error("failed to stop plugin process", "error", err, "plugin", plugin.PluginID)
			return custom_errors.NewPluginError("failed to stop plugins")
		}
	}
	return nil
}

func (p *PluginManager) StopPluginProcess(plugin *domains.PluginInfo) error {
	binDir := filepath.Dir(plugin.Bin)
	pidFile := filepath.Join(binDir, ".pid")
	if _, err := os.Stat(pidFile); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		PluginLogger.Error("failed to stat pid file", "error", err)
		return custom_errors.NewPluginError("failed to stop plugin process")
	}

	pidBytes, err := os.ReadFile(pidFile)
	if err != nil {
		PluginLogger.Error("failed to read pid file", "error", err)
		return custom_errors.NewPluginError("failed to stop plugin process")
	}

	pidStr := strings.TrimSpace(string(pidBytes))
	if pidStr != "" {
		pid, err := strconv.Atoi(pidStr)
		if err == nil && pid > 0 {
			process, err := os.FindProcess(pid)
			if err == nil {
				if err := process.Kill(); err != nil {
					PluginLogger.Warn("failed to kill process (might be already dead)", "error", err)
				}
			}
		}
	}

	if err := os.Remove(pidFile); err != nil {
		PluginLogger.Error("failed to remove pid file", "error", err)
		return custom_errors.NewPluginError("failed to stop plugin process")
	}

	return nil
}

func (p *PluginManager) HasConnectCode(connectCode string) bool {
	return p.connectCodeMgr.HasConnectCode(connectCode)
}

func (p *PluginManager) GetPluginByConnectCode(connectCode string) *domains.PluginInfo {
	return p.connectCodeMgr.GetPluginFromCode(connectCode)
}

func (p *PluginManager) RemovePlugin(pluginID string) error {
	mapper := &mappers.PluginMapper{}
	plugin, err := p.pluginRepo.GetPluginByID(pluginID)
	if err != nil {
		return custom_errors.NewPluginError("Failed to remove plugin")
	}

	pluginInfo := mapper.ToDomain(plugin)
	PluginLogger.Info("Try stop plugin", "plugin", pluginInfo.PluginID)
	p.StopPluginProcess(pluginInfo)

	err = p.pluginRepo.DeletePluginByID(plugin.ID.String())
	if err != nil {
		return custom_errors.NewPluginError("Failed to remove plugin")
	}

	pluginDirName := p.GetPluginDirName(&nagare_plugin.PluginMetadata{
		ID:      pluginInfo.PluginID,
		Version: pluginInfo.Version,
	})
	pluginDirPath := filepath.Join(paths.PluginDir, pluginDirName)
	err = os.RemoveAll(pluginDirPath)
	if err != nil {
		PluginLogger.Error("Failed to remove plugin data", "plugin", pluginInfo.PluginID, "error", err)
		return custom_errors.NewPluginError("Failed to remove plugin")
	}

	return nil
}

func NewPluginManager(pluginRepo *repositories.PluginRepository, connectCodeMgr *ConnectCodeManager) *PluginManager {
	return &PluginManager{
		pluginRepo:     pluginRepo,
		connectCodeMgr: connectCodeMgr,
	}
}
