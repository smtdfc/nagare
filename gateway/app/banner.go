package app

import (
	"fmt"
	"io"
	"os"
)

var Version = "dev"

func printStartupMessage(w io.Writer, addr string, debugMode bool) {
	mode := "production"
	modeColor := "\x1b[38;2;52;211;153m"

	if debugMode {
		mode = "development"
		modeColor = "\x1b[38;2;251;191;36m"
	}

	const (
		reset  = "\x1b[0m"
		muted  = "\x1b[38;2;120;135;165m"
		cyan   = "\x1b[38;2;34;211;238m"
		blue   = "\x1b[38;2;96;165;250m"
		purple = "\x1b[38;2;192;132;252m"
		green  = "\x1b[38;2;74;222;128m"
		white  = "\x1b[38;2;226;232;240m"
	)

	_, _ = fmt.Fprintf(w, `
  %s███╗   ██╗ █████╗  ██████╗  █████╗ ██████╗ ███████╗%s
  %s████╗  ██║██╔══██╗██╔════╝ ██╔══██╗██╔══██╗██╔════╝%s
  %s██╔██╗ ██║███████║██║  ███╗███████║██████╔╝█████╗  %s
  %s██║╚██╗██║██╔══██║██║   ██║██╔══██║██╔══██╗██╔══╝  %s
  %s██║ ╚████║██║  ██║╚██████╔╝██║  ██║██║  ██║███████╗%s
  %s╚═╝  ╚═══╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝%s

  %sG A T E W A Y%s  %sversion: %s%s

  %s●%s  %sAddress   %shttp://%s
  %s●%s  %sMode      %s%s%s
  %s●%s  %sPID       %s%d%s

  %s◆%s  %sGateway initialized%s

`,
		cyan, reset,
		blue, reset,
		blue, reset,
		purple, reset,
		purple, reset,
		cyan, reset,
		muted, reset, white, Version, reset,
		cyan, reset, muted, white, addr,
		blue, reset, muted, modeColor, mode, reset,
		purple, reset, muted, white, os.Getpid(), reset,
		green, reset, green, reset,
	)
}
