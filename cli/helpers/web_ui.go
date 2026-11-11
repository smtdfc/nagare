package helpers

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/smtdfc/nagare/pkgs/paths"
)

//go:embed web_dist/*
var Files embed.FS

const webDaemonEnv = "NAGARE_WEB_DAEMON"

type spaHandler struct {
	fs fs.FS
}

func (h *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cleanPath := path.Clean(r.URL.Path)
	fsPath := strings.TrimPrefix(cleanPath, "/")
	if fsPath == "" || fsPath == "." {
		fsPath = "index.html"
	}

	f, err := h.fs.Open(fsPath)
	if err == nil {
		defer f.Close()

		stat, err := f.Stat()
		if err == nil && !stat.IsDir() {
			ext := filepath.Ext(fsPath)
			if mimeType := mime.TypeByExtension(ext); mimeType != "" {
				w.Header().Set("Content-Type", mimeType)
			}

			http.ServeContent(w, r, stat.Name(), stat.ModTime(), f.(io.ReadSeeker))
			return
		}
	}

	if path.Ext(fsPath) != "" {
		http.NotFound(w, r)
		return
	}

	indexFile, err := h.fs.Open("index.html")
	if err != nil {
		http.Error(w, "index.html not found", http.StatusInternalServerError)
		return
	}
	defer indexFile.Close()

	stat, _ := indexFile.Stat()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile.(io.ReadSeeker))
}

func Handler() http.Handler {
	subFS, err := fs.Sub(Files, "web_dist")
	if err != nil {
		panic("embedded web UI files are unavailable: " + err.Error())
	}

	return &spaHandler{
		fs: subFS,
	}
}

func StartServer(addr string, onReady func()) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	server := &http.Server{
		Handler: Handler(),
	}

	log.Printf("Nagare Web UI listening on http://%s", addr)
	if onReady != nil {
		onReady()
	}

	return server.Serve(listener)
}

func IsWebDaemon() bool {
	return os.Getenv(webDaemonEnv) == "1"
}

func StartBackgroundWebUI(addr string) error {
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to resolve the CLI executable: %w", err)
	}

	logFile := filepath.Join(paths.LogDir, "web-ui.log")
	if err := os.MkdirAll(paths.LogDir, 0o755); err != nil {
		return fmt.Errorf("failed to create the log directory: %w", err)
	}
	logHandle, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", logFile, err)
	}
	defer logHandle.Close()

	cmd := exec.Command(executable, os.Args[1:]...)
	cmd.Env = append(os.Environ(), webDaemonEnv+"=1")
	cmd.SysProcAttr = detachSysProcAttr()
	cmd.Stdout = logHandle
	cmd.Stderr = logHandle
	cmd.Stdin = nil

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start the Web UI daemon: %w", err)
	}

	if err := saveWebUIPID(cmd.Process.Pid); err != nil {
		_ = cmd.Process.Kill()
		return fmt.Errorf("failed to save Web UI PID: %w", err)
	}

	fmt.Printf("Starting Nagare Web UI in the background with PID %d...\n", cmd.Process.Pid)
	fmt.Printf("Logs: %s\n", logFile)

	if err := waitForWebUI(addr, 3*time.Second); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}

	return nil
}

func waitForWebUI(addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 300*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("the Web UI did not start listening on %s within %s", addr, timeout)
}

func saveWebUIPID(pid int) error {
	return os.WriteFile(paths.WebUIPIDFile, []byte(strconv.Itoa(pid)), 0o644)
}

func removeWebUIPID() error {
	err := os.Remove(paths.WebUIPIDFile)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func TryStopWebUI() (bool, error) {
	data, err := os.ReadFile(paths.WebUIPIDFile)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to read Web UI PID file: %w", err)
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil || pid <= 0 {
		_ = removeWebUIPID()
		return false, fmt.Errorf("invalid Web UI PID in %s", paths.WebUIPIDFile)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		_ = removeWebUIPID()
		return false, nil
	}

	if err := requestWebUIShutdown(process); err != nil {
		// The PID file may outlive the daemon when the machine restarts or
		// the daemon exits through its own signal handler; treat a finished
		// process as "nothing to stop" instead of surfacing an error.
		if isProcessFinished(err) {
			_ = removeWebUIPID()
			return false, nil
		}
		_ = removeWebUIPID()
		return false, fmt.Errorf("failed to stop Web UI process %d: %w", pid, err)
	}
	_ = removeWebUIPID()

	return true, nil
}

func isProcessFinished(err error) bool {
	return errors.Is(err, os.ErrProcessDone) || strings.Contains(err.Error(), "process already finished")
}

func InstallWebDaemonSignalCleanup() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		if _, ok := <-signals; ok {
			_ = removeWebUIPID()
			os.Exit(0)
		}
	}()
}
