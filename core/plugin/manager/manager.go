package manager

import (
	"archive/zip"
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/config"
	"github.com/smtdfc/nagare/core/exceptions"
	"github.com/smtdfc/nagare/core/utils"
	"github.com/smtdfc/nagare/plugin-sdk/plugin"
	"github.com/smtdfc/nagare/plugin-sdk/shared"
	nagare_logger "github.com/smtdfc/nagare/shared/logger"
	nagare_path "github.com/smtdfc/nagare/shared/path"
	"go.uber.org/multierr"
)

func GetCurrentPlatform() string {
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}

type PluginData struct {
	Cmd    *exec.Cmd
	Config *config.Plugin
	ID     string
	Name   string
}

type PluginManager struct {
	mu              *sync.RWMutex
	Conf            *config.Config
	PluginInstances map[string]*PluginData
	Arena           map[string]interface{}
	logger          *slog.Logger
}

func (m *PluginManager) Install(pluginPackPath string) error {
	id := uuid.New()
	platform := GetCurrentPlatform()
	pluginPath := filepath.Join(nagare_path.TempDir, id.String())

	err := extractPlugin(pluginPackPath, pluginPath)
	if err != nil {
		return fmt.Errorf("failed to extract plugin: %w", err)
	}

	pluginMetadataFile := filepath.Join(pluginPath, "metadata.json")
	var metadata plugin.PluginMetadata
	err = utils.ReadJSON(pluginMetadataFile, &metadata)
	if err != nil {
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	bin, ok := metadata.Bin[platform]
	if !ok {
		return fmt.Errorf("plugin not support platform: %s", platform)
	}

	binPath := filepath.Join(pluginPath, bin)
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		return fmt.Errorf("plugin binary not found at: %s", binPath)
	}

	pluginDirName := fmt.Sprintf("%s.%s", metadata.ID, base64.URLEncoding.EncodeToString([]byte(metadata.ID)))
	pluginDest := filepath.Join(nagare_path.PluginDir, pluginDirName)

	if m.Conf.Plugins == nil {
		m.Conf.Plugins = map[string]config.Plugin{}
	}

	m.Conf.Plugins[metadata.ID] = config.Plugin{
		Name:     metadata.Name,
		Id:       metadata.ID,
		Author:   metadata.Author,
		Version:  metadata.Version,
		Path:     pluginDirName,
		Features: metadata.Features,
		Bin:      bin,
	}

	err = utils.CopyDir(pluginPath, pluginDest)
	if err != nil {
		return err
	}

	return config.SaveConfig(m.Conf)
}

func (m *PluginManager) LoadPlugin() error {
	var pluginErr error
	for pluginName, pluginConfig := range m.Conf.Plugins {
		pluginPath := filepath.Join(nagare_path.PluginDir, pluginConfig.Path, pluginConfig.Bin)
		id := uuid.New().String()
		cmd := exec.Command(pluginPath)
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("%s=%s", shared.PLUGIN_ID_ENV, id),
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Start()
		if err != nil {
			pluginErr = multierr.Append(pluginErr, exceptions.NewPluginException(fmt.Sprintf("Cannot start plugin %s: %s", pluginName, err), pluginName))
		}

		m.logger.Info("Started plugin", "plugin", pluginName)
		m.PluginInstances[id] = &PluginData{
			Cmd:    cmd,
			Config: &pluginConfig,
			ID:     id,
			Name:   pluginName,
		}
	}

	return pluginErr
}

func extractPlugin(zipPath, destDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, f := range reader.File {
		fpath := filepath.Join(destDir, f.Name)
		if !filepath.HasPrefix(fpath, filepath.Clean(destDir)+string(filepath.Separator)) {
			return fmt.Errorf("illegal file path: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *PluginManager) Shutdown() {
	// GlobalPluginHost.Shutdown()
	for id, instance := range m.PluginInstances {
		err := instance.Cmd.Process.Kill()
		if err != nil {
			m.logger.Info("Kill plugin failed", "plugin", instance.Name, "instance", id)
		}

		m.logger.Info("Killed plugin", "plugin", instance.Name, "instance", id)
	}
}

func NewPluginManager(conf *config.Config) *PluginManager {
	return &PluginManager{
		mu:              &sync.RWMutex{},
		Conf:            conf,
		PluginInstances: make(map[string]*PluginData),
		Arena:           map[string]interface{}{},
		logger:          nagare_logger.GetLogger("Plugin manager"),
	}
}
