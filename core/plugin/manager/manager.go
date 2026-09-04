package manager

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/mappers"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
	"github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/plugin/metadata"
	"github.com/smtdfc/nagare/shared/helpers"
	"github.com/smtdfc/nagare/shared/paths"
)

func unpackPlugin(archivePath, destDir string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open plugin archive: %w", err)
	}
	defer reader.Close()

	if err := os.MkdirAll(destDir, 0775); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	for _, file := range reader.File {
		cleanDestDir := filepath.Clean(destDir)
		fpath := filepath.Join(cleanDestDir, file.Name)

		if !strings.HasPrefix(fpath, cleanDestDir+string(filepath.Separator)) && fpath != cleanDestDir {
			return fmt.Errorf("illegal file path detected in archive: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, 0775); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0775); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		rc.Close()
		outFile.Close()

		if err != nil {
			return err
		}

		if err := os.Chmod(fpath, file.Mode()); err != nil {
			return err
		}
	}

	return nil
}

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

type PluginManager struct {
	pluginRepo   *repositories.PluginRepository
	pluginMapper *mappers.PluginMapper
	connectCodes map[string]string
	logger       *logger.BaseLogger
}

func (p *PluginManager) GetListPlugin(ctx context.Context) ([]*plugin.Plugin, error) {
	ents, err := p.pluginRepo.GetAllPlugin(ctx)
	if err != nil {
		return nil, custom_errors.ErrGetListPluginFailed
	}

	return p.pluginMapper.ToDomains(ents), nil
}

func (p *PluginManager) StopPlugin(ctx context.Context, plugin *plugin.Plugin) error {
	pidPath := plugin.Bin + ".pid"

	data, err := os.ReadFile(pidPath)
	if err != nil {
		if os.IsNotExist(err) {
			p.logger.Warn("PID file not found, skipping plugin stop", "plugin", plugin.PluginID)
			return nil
		}
		p.logger.Error("Failed to read PID file", "error", err, "plugin", plugin.PluginID)
		return err
	}

	var pid int
	if _, err := fmt.Sscanf(string(data), "%d", &pid); err != nil {
		p.logger.Error("Invalid PID format in file", "error", err, "plugin", plugin.PluginID)
		_ = os.Remove(pidPath)
		return err
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		p.logger.Error("Failed to find process", "pid", pid, "error", err, "plugin", plugin.PluginID)
		_ = os.Remove(pidPath)
		return err
	}

	if err := process.Kill(); err != nil {
		p.logger.Error("Failed to kill process", "pid", pid, "error", err, "plugin", plugin.PluginID)
	}

	if err := os.Remove(pidPath); err != nil {
		p.logger.Warn("Failed to remove PID file", "error", err, "plugin", plugin.PluginID)
	}

	delete(p.connectCodes, plugin.PluginID)
	return nil
}

func (p *PluginManager) StartPlugin(ctx context.Context, plugin *plugin.Plugin) error {
	connectCode := uuid.New().String()
	p.connectCodes[plugin.PluginID] = connectCode
	cmd := exec.Command(plugin.Bin)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Env = append(os.Environ(), fmt.Sprintf("NAGARE_PLUGIN_CONNECT_CODE=%s", connectCode))

	err := cmd.Start()
	if err != nil {
		p.logger.Error("Start plugin failed", "error", err, "plugin", plugin.PluginID)
		return custom_errors.ErrStartPluginFailed
	}

	pidPath := plugin.Bin + ".pid"
	pidStr := fmt.Sprintf("%d", cmd.Process.Pid)
	if writeErr := os.WriteFile(pidPath, []byte(pidStr), 0644); writeErr != nil {
		p.logger.Error("Failed to write plugin PID file", "error", writeErr, "plugin", plugin.PluginID)
	}

	return nil
}
func (p *PluginManager) Install(ctx context.Context, pluginPath string) error {
	_, err := os.Stat(pluginPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return custom_errors.ErrPluginNotFound
		}

		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return custom_errors.ErrInstallPluginFailed
	}

	id := uuid.New().String()
	tempDir := filepath.Join(paths.TempDir, id)
	metadataFile := filepath.Join(tempDir, "metadata.json")
	os.MkdirAll(tempDir, 0775)

	err = unpackPlugin(
		pluginPath,
		tempDir,
	)
	if err != nil {
		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return custom_errors.ErrUnpackPluginFailed
	}

	pluginMetadata, err := loadMetadata(metadataFile)
	if err != nil {
		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return custom_errors.ErrLoadPluginMetadataFailed
	}

	err = pluginMetadata.Validate()
	if err != nil {
		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return custom_errors.ErrPluginMetadataInvalid
	}

	pluginDir := filepath.Join(paths.PluginDir, pluginMetadata.ID)
	binFile, isExist := pluginMetadata.Bin[runtime.GOOS]
	if !isExist {
		return custom_errors.ErrPluginBinaryMissing
	}

	err = helpers.CopyDir(tempDir, pluginDir)
	if err != nil {
		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return custom_errors.ErrInstallPluginFailed
	}

	binFile = filepath.Join(pluginDir, binFile)
	if err := os.Chmod(binFile, 0755); err != nil {
		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return custom_errors.ErrInstallPluginFailed
	}

	plugin := plugin.Plugin{
		PluginID: pluginMetadata.ID,
		Name:     pluginMetadata.Name,
		Author:   pluginMetadata.Author,
		Features: []plugin.PluginFeature{},
		Version:  pluginMetadata.Version,
		Bin:      binFile,
		IsActive: true,
	}

	_, err = p.pluginRepo.CreateOrUpdate(ctx, p.pluginMapper.ToEntity(&plugin))
	if err != nil {
		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return custom_errors.ErrInstallPluginFailed
	}

	return p.StartPlugin(ctx, &plugin)
}

func (p *PluginManager) StartAllPlugin(ctx context.Context) error {
	activePlugins, err := p.pluginRepo.GetAllActivePlugin(ctx)
	if err != nil {
		p.logger.Error("Start plugin failed", "error", err)
		return custom_errors.ErrStartPluginFailed
	}

	for _, plugin := range p.pluginMapper.ToDomains(activePlugins) {
		p.StartPlugin(ctx, plugin)
	}

	return nil
}

func (p *PluginManager) StopAllPlugin(ctx context.Context) error {
	activePlugins, err := p.pluginRepo.GetAllActivePlugin(ctx)
	if err != nil {
		p.logger.Error("Stop plugin failed", "error", err)
		return custom_errors.ErrStartPluginFailed
	}

	for _, plugin := range p.pluginMapper.ToDomains(activePlugins) {
		p.StopPlugin(ctx, plugin)
	}

	return nil
}

// @Injectable
func NewPluginManager(
	pluginRepo *repositories.PluginRepository,
	pluginMapper *mappers.PluginMapper,
	logger *logger.BaseLogger,
) *PluginManager {
	return &PluginManager{
		pluginRepo:   pluginRepo,
		logger:       logger,
		connectCodes: make(map[string]string),
	}
}
