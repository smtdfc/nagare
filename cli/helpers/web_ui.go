package helpers

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/smtdfc/nagare/pkgs/paths"
)

//go:embed web_dist/*
var Files embed.FS

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

// StartServer serves the Web UI on addr. When the listener is ready, onReady
// runs before the server blocks serving requests, so callers can open a browser
// without racing the bind.
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

// webDaemonEnv marks the detached child process so the web command serves
// requests instead of spawning another background process.
const webDaemonEnv = "NAGARE_WEB_DAEMON"

// IsWebDaemon reports whether the current process is the detached Web UI
// daemon started by StartBackgroundWebUI.
func IsWebDaemon() bool {
	return os.Getenv(webDaemonEnv) == "1"
}

// StartBackgroundWebUI re-executes the nagare CLI as a detached daemon with
// the same arguments, redirecting its output to a log file. The parent
// process returns as soon as the Web UI answers on its address.
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

	fmt.Printf("Starting Nagare Web UI in the background with PID %d...\n", cmd.Process.Pid)
	fmt.Printf("Logs: %s\n", logFile)

	if err := waitForWebUI(addr, 3*time.Second); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}

	return nil
}

// waitForWebUI polls the Web UI address until it accepts TCP connections,
// so the CLI can report a failure to bind instead of exiting silently.
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
