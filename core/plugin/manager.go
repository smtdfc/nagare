package plugin

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/agent"
	"github.com/smtdfc/nagare/core/config"
	"github.com/smtdfc/nagare/core/utils"
	"github.com/smtdfc/nagare/plugin-sdk/host"
	"github.com/smtdfc/nagare/plugin-sdk/plugin"
	"github.com/smtdfc/nagare/plugin-sdk/shared"
	nagare_logger "github.com/smtdfc/nagare/shared/logger"
	nagare_path "github.com/smtdfc/nagare/shared/path"
	"go.uber.org/multierr"
)

var ChatMgr *ChatChannelManager

func GetCurrentPlatform() string {
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}

type PluginManager struct {
	mu              *sync.RWMutex
	Conf            *config.Config
	PluginInstances map[string]*exec.Cmd
	Host            *host.Host
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

func (m *PluginManager) StartHost(pool *agent.AgentPool, sessionMgr *agent.SessionManager) {
	m.Host.Handler(shared.REGISTER_CHAT_CHANNEL, func(msg shared.Message) {
		if ChatMgr == nil {
			ChatMgr = NewChatChannelManager(m.Host)
		}

		var payload shared.RegisterChatChannelPayload
		json.Unmarshal(msg.Payload, &payload)
		a := pool.GetOrNew()
		a.State.ExtendHistory(
			sessionMgr.GetHistory(payload.ID, agent.NAGARE_LIST_MESSAGE_SIZE_LIMIT),
		)
		ChatMgr.Register(&ChatChannel{
			Id:         payload.ID,
			Agent:      a,
			SessionMgr: sessionMgr,
			SessionID:  payload.ID,
			CleanUp: func() {
				sessionMgr.SaveHistory(payload.ID, a.State.History)
				pool.Put(a)
			},
		})

		m.Host.Send(msg.PluginID, shared.REGISTER_CHAT_CHANNEL_SUCCESS, shared.RegisterChatChannelSuccessPayload{
			ID: payload.ID,
		})
	})

	m.Host.Handler(shared.HANDLE_CHAT_MESSAGE, func(msg shared.Message) {
		var payload shared.HandleChatMessagePayload
		json.Unmarshal(msg.Payload, &payload)
		ChatMgr.Handle(&payload, msg.PluginID)
	})

	m.Host.Start()
}

func (m *PluginManager) LoadPlugin() error {
	var pluginErr error
	for pluginName, pluginConfig := range m.Conf.Plugins {
		pluginPath := filepath.Join(nagare_path.PluginDir, pluginConfig.Path, pluginConfig.Bin)
		cmd := exec.Command(pluginPath)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Start()
		if err != nil {
			pluginErr = multierr.Append(pluginErr, fmt.Errorf("Cannot start plugin %s: %s", pluginName, err))
		}
		m.logger.Info("Started plugin", "plugin", pluginName)
		m.PluginInstances[pluginName] = cmd
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
	m.Host.Shutdown()
	for name, instance := range m.PluginInstances {
		err := instance.Process.Kill()
		if err != nil {
			m.logger.Info("Kill plugin failed", "plugin", name)
		}

		m.logger.Info("Killed plugin", "plugin", name)
	}
}

func NewPluginManager(conf *config.Config) *PluginManager {
	return &PluginManager{
		mu:              &sync.RWMutex{},
		Conf:            conf,
		PluginInstances: map[string]*exec.Cmd{},
		Arena:           map[string]interface{}{},
		Host:            host.NewHost(nagare_logger.GetLogger("Plugin Host")),
		logger:          nagare_logger.GetLogger("Plugin manager"),
	}
}
