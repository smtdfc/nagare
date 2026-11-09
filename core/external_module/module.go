package external_module

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/smtdfc/nagare/pkgs/paths"
)

const (
	moduleIDEnv          = "NAGARE_MODULE_ID"
	moduleConnectCodeEnv = "NAGARE_MODULE_CONNECT_CODE"
	moduleSocketPathEnv  = "NAGARE_MODULE_SOCKET_PATH"
	moduleDataDirEnv     = "NAGARE_MODULE_DATA_DIR"
	moduleConfigDirEnv   = "NAGARE_MODULE_CONFIG_DIR"
	moduleLogFileEnv     = "NAGARE_MODULE_LOG_FILE"
)

func stopGracefully(p *os.Process) {
	if p == nil || !alive(p) {
		return
	}
	if err := terminate(p); err != nil {
		return
	}

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if !alive(p) {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	_ = p.Kill()
}

type ExternalModule struct {
	ID        string
	Bin       string
	DataDir   string
	ConfigDir string
}

func (e *ExternalModule) Start(connectCode string) error {
	if _, err := os.Stat(e.Bin); err != nil {
		return fmt.Errorf("external module %s: binary %s is not available: %w", e.ID, e.Bin, err)
	}

	for _, dir := range []string{e.DataDir, e.ConfigDir, paths.ModuleLogDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("external module %s: failed to create directory %s: %w", e.ID, dir, err)
		}
	}

	cmd := exec.Command(e.Bin)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(
		os.Environ(),
		fmt.Sprintf("%s=%s", moduleIDEnv, e.ID),
		fmt.Sprintf("%s=%s", moduleConnectCodeEnv, connectCode),
		fmt.Sprintf("%s=%s", moduleSocketPathEnv, paths.GetModuleSocketPath(e.ID)),
		fmt.Sprintf("%s=%s", moduleDataDirEnv, e.DataDir),
		fmt.Sprintf("%s=%s", moduleConfigDirEnv, e.ConfigDir),
		fmt.Sprintf("%s=%s", moduleLogFileEnv, filepath.Join(paths.ModuleLogDir, e.ID+".log")),
	)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("external module %s: failed to start %s: %w", e.ID, e.Bin, err)
	}

	go func() {
		_ = cmd.Wait()
	}()

	pidPath := e.pidFile()
	if err := os.WriteFile(pidPath, []byte(strconv.Itoa(cmd.Process.Pid)), 0644); err != nil {
		// The process is already running; do not leak it just because the PID
		// could not be recorded.
		_ = cmd.Process.Kill()
		return fmt.Errorf("external module %s: failed to write pid file %s: %w", e.ID, pidPath, err)
	}

	return nil
}

func (e *ExternalModule) Stop() error {
	pidPath := e.pidFile()

	data, err := os.ReadFile(pidPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("external module %s: failed to read pid file %s: %w", e.ID, pidPath, err)
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		_ = os.Remove(pidPath)
		return fmt.Errorf("external module %s: invalid pid file %s: %w", e.ID, pidPath, err)
	}

	proc, err := os.FindProcess(pid)
	if err == nil {
		stopGracefully(proc)
	}

	if err := os.Remove(pidPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("external module %s: failed to remove pid file %s: %w", e.ID, pidPath, err)
	}

	return nil
}

func (e *ExternalModule) pidFile() string {
	return e.Bin + ".pid"
}

func NewExternalModule(id string) ExternalModule {
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	return ExternalModule{
		ID:        id,
		Bin:       filepath.Join(paths.ModuleDir, id+ext),
		DataDir:   filepath.Join(paths.ModuleDataDir, id),
		ConfigDir: filepath.Join(paths.ModuleConfigDir, id),
	}
}
