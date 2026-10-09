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
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/shirou/gopsutil/v4/process"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/mappers"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
	"github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/pkgs/helpers"
	"github.com/smtdfc/nagare/pkgs/paths"
	"github.com/smtdfc/nagare/plugin/metadata"
)

func extractFile(fpath string, file *zip.File) error {
	outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {

		}
	}(outFile)

	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer func(rc io.ReadCloser) {
		err := rc.Close()
		if err != nil {

		}
	}(rc)

	if _, err = io.Copy(outFile, rc); err != nil {
		return err
	}

	return os.Chmod(fpath, file.Mode())
}

func unpackPlugin(archivePath, destDir string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open plugin archive: %w", err)
	}
	defer func(reader *zip.ReadCloser) {
		_ = reader.Close()
	}(reader)

	absDestDir, err := filepath.Abs(destDir)
	if err != nil {
		return fmt.Errorf("failed to resolve destination directory: %w", err)
	}
	absDestDir = filepath.Clean(absDestDir)
	if err := os.MkdirAll(absDestDir, 0775); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	for _, file := range reader.File {
		if filepath.IsAbs(file.Name) {
			return fmt.Errorf("illegal absolute file path detected in archive: %s", file.Name)
		}

		normalizedName := strings.ReplaceAll(file.Name, "\\", "/")
		if strings.Contains(normalizedName, "../") || normalizedName == ".." {
			return fmt.Errorf("illegal file path detected in archive: %s", file.Name)
		}

		fpath := filepath.Clean(filepath.Join(absDestDir, file.Name))
		destPrefix := absDestDir + string(filepath.Separator)
		if fpath != absDestDir && !strings.HasPrefix(fpath, destPrefix) {
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

		if err := extractFile(fpath, file); err != nil {
			return err
		}
	}

	return nil
}

func loadMetadata(metadataFile string) (*metadata.PluginMetadata, error) {
	_, err := os.Stat(metadataFile)
	if errors.Is(err, os.ErrNotExist) {
		return nil, errors.New("metadata file not exist")
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
	socketPath   string
	logger       *logger.BaseLogger
}

func (p *PluginManager) getPluginByID(ctx context.Context, id string) (*plugin.Plugin, error) {
	pluginEntity, err := p.pluginRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if pluginEntity == nil {
		return nil, custom_errors.ErrPluginNotFound
	}

	return p.pluginMapper.ToDomain(pluginEntity), nil
}

func (p *PluginManager) GetListPlugin(ctx context.Context) ([]*plugin.Plugin, error) {
	ents, err := p.pluginRepo.FindAll(ctx)
	if err != nil {
		return nil, custom_errors.ErrGetListPluginFailed
	}

	return p.pluginMapper.ToDomains(ents), nil
}

func (p *PluginManager) StopPlugin(plugin *plugin.Plugin) error {
	pidPath := plugin.Bin + ".pid"

	data, err := os.ReadFile(pidPath)
	if err != nil {
		if os.IsNotExist(err) {
			p.logger.Warn("PID file not found, skipping plugin stop", "packageName", plugin.PackageName)
			return nil
		}
		p.logger.Error("Failed to read PID file", "error", err, "packageName", plugin.PackageName)
		return err
	}

	var pid int
	if _, err := fmt.Sscanf(string(data), "%d", &pid); err != nil {
		p.logger.Error("Invalid PID format in file", "error", err, "packageName", plugin.PackageName)
		_ = os.Remove(pidPath)
		return err
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		p.logger.Error("Failed to find process", "pid", pid, "error", err, "packageName", plugin.PackageName)
		_ = os.Remove(pidPath)
		return err
	}

	if err := proc.Kill(); err != nil {
		p.logger.Error("Failed to kill process", "pid", pid, "error", err, "packageName", plugin.PackageName)
	}

	if err := os.Remove(pidPath); err != nil {
		p.logger.Warn("Failed to remove PID file", "error", err, "packageName", plugin.PackageName)
	}

	delete(p.connectCodes, plugin.PackageName)
	return nil
}

func (p *PluginManager) StartPlugin(plugin *plugin.Plugin) error {
	connectCode := uuid.New().String()
	p.connectCodes[plugin.PackageName] = connectCode
	cmd := exec.Command(plugin.Bin)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	socketPath := p.socketPath
	if socketPath == "" {
		socketPath = paths.PluginSocketPath
	}
	cmd.Env = append(
		os.Environ(),
		fmt.Sprintf("NAGARE_PLUGIN_CONNECT_CODE=%s", connectCode),
		fmt.Sprintf("NAGARE_PLUGIN_SOCKET_PATH=%s", socketPath),
		fmt.Sprintf("NAGARE_PLUGIN_LOG_FILE=%s", filepath.Join(paths.PluginLogDir, plugin.PackageName+".log")),
		fmt.Sprintf("NAGARE_PLUGIN_CONFIG_DIR=%s", filepath.Join(paths.PluginConfigDir, plugin.PackageName)),
	)

	err := cmd.Start()
	if err != nil {
		p.logger.Error("Start plugin failed", "error", err, "packageName", plugin.PackageName)
		return custom_errors.ErrStartPluginFailed
	}

	pidPath := plugin.Bin + ".pid"
	pidStr := fmt.Sprintf("%d", cmd.Process.Pid)
	if writeErr := os.WriteFile(pidPath, []byte(pidStr), 0644); writeErr != nil {
		p.logger.Error("Failed to write plugin PID file", "error", writeErr, "packageName", plugin.PackageName)
	}

	return nil
}

func (p *PluginManager) Install(ctx context.Context, pluginPath string) (*plugin.Plugin, error) {
	_, err := os.Stat(pluginPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, custom_errors.ErrPluginNotFound
		}

		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrInstallPluginFailed
	}

	id := uuid.New().String()
	tempDir := filepath.Join(paths.TempDir, id)
	metadataFile := filepath.Join(tempDir, "metadata.json")
	err = os.MkdirAll(tempDir, 0775)
	if err != nil {
		return nil, err
	}

	err = unpackPlugin(
		pluginPath,
		tempDir,
	)
	if err != nil {
		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrUnpackPluginFailed
	}

	pluginMetadata, err := loadMetadata(metadataFile)
	if err != nil {
		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrLoadPluginMetadataFailed
	}

	err = pluginMetadata.Validate()
	if err != nil {
		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrPluginMetadataInvalid
	}

	pluginDir := filepath.Join(paths.PluginDir, pluginMetadata.PackageName)
	binFile, isExist := pluginMetadata.Bin[runtime.GOOS]
	if !isExist {
		return nil, custom_errors.ErrPluginBinaryMissing
	}

	err = helpers.CopyDir(tempDir, pluginDir)
	if err != nil {
		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrInstallPluginFailed
	}

	binFile = filepath.Join(pluginDir, binFile)
	if err := os.Chmod(binFile, 0755); err != nil {
		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrInstallPluginFailed
	}

	fmt.Println("features", pluginMetadata.Features)
	newPlugin := plugin.Plugin{
		PackageName: pluginMetadata.PackageName,
		Name:        pluginMetadata.Name,
		Author:      pluginMetadata.Author,
		Features:    plugin.ParseFeatureString(pluginMetadata.GetFeatures()),
		Version:     pluginMetadata.Version,
		Bin:         binFile,
		IsActive:    true,
	}

	_, err = p.pluginRepo.CreateOrUpdate(ctx, p.pluginMapper.ToEntity(&newPlugin))
	if err != nil {
		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrInstallPluginFailed
	}

	err = p.StartPlugin(&newPlugin)
	if err != nil {
		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrStartPluginFailed
	}

	return &newPlugin, nil

}

func (p *PluginManager) StartAllPlugin(ctx context.Context) error {
	activePlugins, err := p.pluginRepo.FindActive(ctx)
	if err != nil {
		p.logger.Error("Start plugin failed", "error", err)
		return custom_errors.ErrStartPluginFailed
	}

	for _, plg := range p.pluginMapper.ToDomains(activePlugins) {
		err := p.StartPlugin(plg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PluginManager) SetPluginSocketPath(path string) {
	p.socketPath = path
}

func (p *PluginManager) GetPluginSocketPath() string {
	return p.socketPath
}

func (p *PluginManager) SetConnectCode(packageName string, code string) {
	p.connectCodes[packageName] = code
}

func (p *PluginManager) StopAllPlugin(ctx context.Context) error {
	activePlugins, err := p.pluginRepo.FindActive(ctx)
	if err != nil {
		p.logger.Error("Stop plugin failed", "error", err)
		return custom_errors.ErrStartPluginFailed
	}

	for _, plg := range p.pluginMapper.ToDomains(activePlugins) {
		err := p.StopPlugin(plg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PluginManager) ValidConnect(cxt context.Context, packageName string, connectCode string) (*plugin.Plugin, error) {
	pluginEntity, err := p.pluginRepo.FindByPackageName(cxt, packageName)
	if err != nil {
		return nil, custom_errors.ErrCheckPluginConnectionFailed
	}

	if pluginEntity == nil {
		return nil, custom_errors.ErrPluginNotFound
	}

	code, ok := p.connectCodes[packageName]
	if !ok {
		return nil, custom_errors.ErrPluginConnectionInvalid
	}

	if code != connectCode {
		return nil, custom_errors.ErrPluginConnectionInvalid
	}

	return p.pluginMapper.ToDomain(pluginEntity), nil
}

func (p *PluginManager) Uninstall(ctx context.Context, id string) error {
	plg, err := p.getPluginByID(ctx, id)
	if err != nil {
		if errors.Is(err, custom_errors.ErrPluginNotFound) {
			return err
		}

		return custom_errors.ErrUninstallPluginFailed
	}

	p.logger.Info("Uninstalling plugin", "packageName", plg.PackageName, "name", plg.Name, "version", plg.Version)

	p.logger.Info("Stopping plugin before uninstall", "packageName", plg.PackageName)
	err = p.StopPlugin(plg)
	if err != nil {
		p.logger.Logger.Error("Failed to stop plugin before uninstall", "error", err, "packageName", plg.PackageName)
	} else {
		p.logger.Info("Plugin stopped successfully", "packageName", plg.PackageName)
	}

	err = p.pluginRepo.DeleteById(ctx, plg.ID.String())
	if err != nil {
		return custom_errors.ErrUninstallPluginFailed
	}

	pluginDir := filepath.Join(paths.PluginDir, plg.PackageName)
	err = os.RemoveAll(pluginDir)
	if err != nil {
		p.logger.Error("Failed to remove plugin directory", "error", err, "packageName", plg.PackageName)
	}

	p.logger.Info("Uninstall plugin completed", "packageName", plg.PackageName, "name", plg.Name, "version", plg.Version)
	return nil
}

func (p *PluginManager) Activate(ctx context.Context, id string) error {
	plg, err := p.getPluginByID(ctx, id)
	if err != nil {
		if errors.Is(err, custom_errors.ErrPluginNotFound) {
			return err
		}

		return custom_errors.ErrActivatePluginFailed
	}

	p.logger.Info("Activating plugin", "packageName", plg.PackageName, "name", plg.Name, "version", plg.Version)

	if plg.IsActive {
		return custom_errors.ErrPluginAlreadyActive
	}

	plg.IsActive = true
	err = p.pluginRepo.Update(ctx, p.pluginMapper.ToEntity(plg))
	if err != nil {
		return custom_errors.ErrActivatePluginFailed
	}

	p.logger.Info("Plugin activated successfully", "packageName", plg.PackageName, "name", plg.Name, "version", plg.Version)
	return p.StartPlugin(plg)
}

func (p *PluginManager) Deactivate(ctx context.Context, id string) error {
	plg, err := p.getPluginByID(ctx, id)
	if err != nil {
		if errors.Is(err, custom_errors.ErrPluginNotFound) {
			return err
		}

		return custom_errors.ErrDeactivatePluginFailed
	}

	p.logger.Info("Deactivating plugin", "packageName", plg.PackageName, "name", plg.Name, "version", plg.Version)
	if !plg.IsActive {
		return custom_errors.ErrPluginNotActive
	}

	err = p.StopPlugin(plg)
	if err != nil {
		return custom_errors.ErrDeactivatePluginFailed
	}

	plg.IsActive = false
	err = p.pluginRepo.Update(ctx, p.pluginMapper.ToEntity(plg))
	if err != nil {
		return custom_errors.ErrDeactivatePluginFailed
	}
	p.logger.Info("Plugin deactivated successfully", "packageName", plg.PackageName, "name", plg.Name, "version", plg.Version)
	return nil
}

func (p *PluginManager) GetPluginStatus(ctx context.Context, id string) (*plugin.Status, error) {
	pluginDomain, err := p.getPluginByID(ctx, id)
	if err != nil {
		if errors.Is(err, custom_errors.ErrPluginNotFound) {
			return nil, err
		}

		return nil, custom_errors.ErrGetPluginStatusFailed
	}

	pidPath := pluginDomain.Bin + ".pid"
	data, err := os.ReadFile(pidPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &plugin.Status{
				PID:         "",
				PackageName: pluginDomain.PackageName,
				Name:        pluginDomain.Name,
				Version:     pluginDomain.Version,
			}, nil
		}
		p.logger.Error("Failed to read PID file", "error", err, "packageName", pluginDomain.PackageName)
		return nil, custom_errors.ErrGetPluginStatusFailed
	}

	var pid string
	if _, err := fmt.Sscanf(string(data), "%s", &pid); err != nil {
		p.logger.Error("Invalid PID format in file", "error", err, "packageName", pluginDomain.PackageName)
		return nil, custom_errors.ErrGetPluginStatusFailed
	}

	pidValue, err := strconv.ParseInt(pid, 10, 32)
	if err != nil {
		p.logger.Error("Invalid PID format in file", "error", err, "packageName", pluginDomain.PackageName)
		return nil, custom_errors.ErrGetPluginStatusFailed
	}

	proc, err := process.NewProcess(int32(pidValue))
	if err != nil {
		p.logger.Error("Failed to get process", "error", err, "packageName", pluginDomain.PackageName)
		return nil, custom_errors.ErrGetPluginStatusFailed
	}
	cpuPercent, _ := proc.CPUPercent()
	memInfo, _ := proc.MemoryInfo()

	return &plugin.Status{
		PID:         pid,
		PackageName: pluginDomain.PackageName,
		Name:        pluginDomain.Name,
		Version:     pluginDomain.Version,
		CPUPercent:  cpuPercent,
		MemoryUsage: float64(memInfo.RSS) / (1024 * 1024),
	}, nil
}

// @Injectable
func NewPluginManager(
	pluginRepo *repositories.PluginRepository,
	pluginMapper *mappers.PluginMapper,
	logger *logger.BaseLogger,
) *PluginManager {
	return &PluginManager{
		pluginRepo:   pluginRepo,
		pluginMapper: pluginMapper,
		logger:       logger.With("module", "plugin-manager"),
		connectCodes: make(map[string]string),
		socketPath:   paths.PluginSocketPath,
	}
}
