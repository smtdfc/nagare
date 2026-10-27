package app

import (
	"fmt"
	"io"
	"os"
)

// Version is set at build time via
// -ldflags "-X github.com/smtdfc/nagare/gateway/app.Version=x.y.z".
var Version = "dev"

// printStartupMessage writes a compact startup summary. It intentionally
// replaces Fiber's ASCII banner, which also dumped internal framework details
// (handler count, prefork state, process count) that are not useful to users.
func printStartupMessage(w io.Writer, addr string, debugMode bool) {
	mode := "production"
	if debugMode {
		mode = "development"
	}

	_, _ = fmt.Fprintf(
		w,
		"\n  nagare gateway %s\n"+
			"  listening   http://%s\n"+
			"  mode        %s\n"+
			"  pid         %d\n\n",
		Version,
		addr,
		mode,
		os.Getpid(),
	)
}
