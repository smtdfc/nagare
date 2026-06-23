package plugin

import (
	"archive/zip"
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/hashicorp/go-plugin"
	"github.com/smtdfc/nagare/core/config"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/utils"
	plugin_sdk "github.com/smtdfc/nagare/plugin-sdk"
	"github.com/smtdfc/nagare/plugin-sdk/proto"
	"go.uber.org/multierr"
)

type PluginManager struct {
	Conf     *config.Config
	Clients  map[string]*plugin.Client
	Services map[string]proto.PluginClient
	Bridges  map[string]interface{}
	logger   *slog.Logger
}

func (m *PluginManager) Install(pluginPackPath string) error {
	id := uuid.New()
	pluginPath := filepath.Join(utils.TempDir, id.String())

	err := extractPlugin(pluginPackPath, pluginPath)
	if err != nil {
		return fmt.Errorf("failed to extract plugin: %w", err)
	}

	pluginMetadataFile := filepath.Join(pluginPath, "metadata.json")
	var metadata plugin_sdk.PluginMetadata
	err = utils.ReadJSON(pluginMetadataFile, &metadata)
	if err != nil {
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	binPath := filepath.Join(pluginPath, metadata.Bin)
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		return fmt.Errorf("plugin binary not found at: %s", binPath)
	}

	pluginDirName := fmt.Sprintf("%s.%s", metadata.ID, base64.URLEncoding.EncodeToString([]byte(metadata.ID)))
	pluginDest := filepath.Join(utils.PluginDir, pluginDirName)

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
		Bin:      metadata.Bin,
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
		pluginPath := filepath.Join(utils.PluginDir, pluginConfig.Path, pluginConfig.Bin)

		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: plugin.HandshakeConfig{
				ProtocolVersion:  plugin_sdk.PLUGIN_PROTOCOL_VERSION,
				MagicCookieKey:   plugin_sdk.PLUGIN_MAGIC_COOKIE_KEY,
				MagicCookieValue: plugin_sdk.PLUGIN_MAGIC_COOKIE_VALUE,
			},
			Plugins: map[string]plugin.Plugin{
				"plugin": &plugin_sdk.NagarePlugin{},
			},
			Cmd:              exec.Command(pluginPath),
			Managed:          true,
			AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		})

		rpcClient, err := client.Client()
		if err != nil {
			client.Kill()
			pluginErr = multierr.Append(pluginErr, fmt.Errorf("failed to get rpc client for %s: %w", pluginName, err))
			continue
		}

		raw, err := rpcClient.Dispense("plugin")
		if err != nil {
			client.Kill()
			pluginErr = multierr.Append(pluginErr, fmt.Errorf("failed to resolve plugin %s: %w", pluginName, err))
			continue
		}

		svc, ok := raw.(proto.PluginClient)
		if !ok {
			client.Kill()
			pluginErr = multierr.Append(pluginErr, fmt.Errorf("plugin %s does not implement PluginClient", pluginName))
			continue
		}

		m.Clients[pluginName] = client
		m.Bridges[pluginName] = raw
		m.Services[pluginName] = svc
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
	for name, client := range m.Clients {
		m.logger.Info("Killing plugin", "plugin", name)
		client.Kill()
	}
}

func NewPluginManager(conf *config.Config) *PluginManager {
	return &PluginManager{
		Conf:     conf,
		Clients:  map[string]*plugin.Client{},
		Services: map[string]proto.PluginClient{},
		Bridges:  map[string]interface{}{},
		logger:   logger.GetLogger("Plugin manager"),
	}
}
