package cmd

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/smtdfc/nagare/cli/helpers"
	"github.com/smtdfc/nagare/pkgs/paths"
	"github.com/spf13/cobra"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorBold   = "\033[1m"
)

type checkResult struct {
	name   string
	ok     bool
	warn   bool
	detail string
}

func printCheck(r checkResult) {
	icon := colorGreen + "✔" + colorReset
	label := colorGreen + "OK" + colorReset
	if r.warn {
		icon = colorYellow + "⚠" + colorReset
		label = colorYellow + "WARN" + colorReset
	}
	if !r.ok {
		icon = colorRed + "✘" + colorReset
		label = colorRed + "FAIL" + colorReset
	}
	fmt.Printf("  %s  %-44s [%s]\n", icon, r.name, label)
	if r.detail != "" {
		fmt.Printf("       %s\n", r.detail)
	}
}

func requiredModules() []string {
	switch runtime.GOOS {
	case "darwin":
		return []string{"libnagare_vector.dylib"}
	case "windows":
		return []string{"nagare_vector.dll"}
	default:
		return []string{"libnagare_vector.so"}
	}
}

func gatewayPort() string {
	if p := os.Getenv("NAGARE_GATEWAY_PORT"); p != "" {
		return p
	}
	return "9832"
}

func checkGatewayProcess() checkResult {
	r := checkResult{name: "Gateway process"}

	data, err := os.ReadFile(paths.GatewayPIDFile)
	if os.IsNotExist(err) {
		r.ok = false
		r.detail = "not running — no PID file found"
		return r
	}
	if err != nil {
		r.ok = false
		r.detail = fmt.Sprintf("cannot read PID file: %v", err)
		return r
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil || pid <= 0 {
		r.ok = false
		r.detail = fmt.Sprintf("invalid PID in %s", paths.GatewayPIDFile)
		return r
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		r.ok = false
		r.detail = fmt.Sprintf("process PID %d not found", pid)
		return r
	}

	if err := process.Signal(syscall.Signal(0)); err != nil {
		r.ok = false
		r.detail = fmt.Sprintf("process PID %d is not alive: %v", pid, err)
		return r
	}

	r.ok = true
	r.detail = fmt.Sprintf("running (PID %d)", pid)
	return r
}

func checkGatewayHTTP() checkResult {
	r := checkResult{name: "Gateway HTTP reachable"}
	port := gatewayPort()
	url := fmt.Sprintf("http://localhost:%s/ws", port)

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		r.ok = false
		r.detail = fmt.Sprintf("cannot reach %s — %v", url, err)
		return r
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// Any response (incl. 426 Upgrade Required) means the gateway is serving.
	r.ok = true
	r.detail = fmt.Sprintf("localhost:%s → HTTP %d", port, resp.StatusCode)
	return r
}

func checkPluginSocket() checkResult {
	r := checkResult{name: "Plugin IPC socket"}
	socketPath := paths.PluginSocketPath

	if _, err := os.Stat(socketPath); err != nil {
		r.ok = false
		r.detail = fmt.Sprintf("socket not found: %s", socketPath)
		return r
	}

	conn, err := net.DialTimeout("unix", socketPath, time.Second)
	if err != nil {
		// Socket file exists but no one is listening — stale socket.
		r.ok = false
		r.warn = false
		r.detail = fmt.Sprintf("socket exists but not connectable: %v", err)
		return r
	}
	err = conn.Close()
	if err != nil {
		return checkResult{}
	}

	r.ok = true
	r.detail = socketPath
	return r
}

var sqliteMagic = []byte("SQLite format 3\x00")

func checkDatabase() checkResult {
	r := checkResult{name: "SQLite database"}
	dbPath := paths.DatabaseDir + "/nagare.db"

	f, err := os.Open(dbPath)
	if err != nil {
		r.ok = false
		r.detail = fmt.Sprintf("cannot open database: %v", err)
		return r
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	header := make([]byte, 16)
	if _, err := f.Read(header); err != nil {
		r.ok = false
		r.detail = fmt.Sprintf("cannot read database header: %v", err)
		return r
	}

	for i, b := range sqliteMagic {
		if header[i] != b {
			r.ok = false
			r.detail = "file exists but does not appear to be a valid SQLite database"
			return r
		}
	}

	r.ok = true
	r.detail = dbPath
	return r
}

func checkKeyring() checkResult {
	r := checkResult{name: "RSA keyring"}
	_, _, err := helpers.GetRSAKey()
	if err != nil {
		r.ok = false
		r.detail = fmt.Sprintf("error: %v", err)
		return r
	}
	r.ok = true
	r.detail = fmt.Sprintf("service=%s", helpers.SERVICE_NAME)
	return r
}

func checkModules() []checkResult {
	var results []checkResult
	for _, name := range requiredModules() {
		path := paths.ModulesDir + "/" + name
		r := checkResult{name: fmt.Sprintf("Module: %s", name)}
		if _, err := os.Stat(path); err == nil {
			r.ok = true
			r.detail = path
		} else {
			r.ok = false
			r.detail = fmt.Sprintf("missing: %s", path)
		}
		results = append(results, r)
	}
	return results
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check Nagare environment and dependencies",
	Long:  "Verify that the Nagare gateway is running, reachable, and all required components are healthy.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("\n%sNagare Doctor%s\n\n", colorBold, colorReset)

		var checks []checkResult
		checks = append(checks, checkGatewayProcess())
		checks = append(checks, checkGatewayHTTP())
		checks = append(checks, checkPluginSocket())
		checks = append(checks, checkDatabase())
		checks = append(checks, checkKeyring())
		checks = append(checks, checkModules()...)

		passed, warned, failed := 0, 0, 0
		for _, c := range checks {
			printCheck(c)
			switch {
			case !c.ok:
				failed++
			case c.warn:
				warned++
			default:
				passed++
			}
		}

		fmt.Printf("  %s%d passed%s  •  %s%d warnings%s  •  %s%d failed%s\n\n",
			colorGreen, passed, colorReset,
			colorYellow, warned, colorReset,
			colorRed, failed, colorReset,
		)

		if failed > 0 {
			fmt.Printf("%sSome checks failed.%s\n\n", colorRed, colorReset)
			os.Exit(1)
		}
		if warned > 0 {
			fmt.Printf("%sAll critical checks passed with warnings.%s\n\n", colorYellow, colorReset)
			return
		}
		fmt.Printf("%sAll checks passed. Nagare is ready.%s\n\n", colorGreen, colorReset)
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
