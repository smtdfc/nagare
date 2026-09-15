package builder

import "fmt"

// ANSI Color Codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
	ColorBold   = "\033[1m"
)

func printStep(msg string) {
	fmt.Printf("%s[*]%s %s ... ", ColorCyan, ColorReset, msg)
}

func printSuccess() {
	fmt.Printf("%s[OK]%s\n", ColorGreen, ColorReset)
}

func printError(format string, a ...interface{}) {
	fmt.Printf("\n%s[ERROR]%s "+format+"\n", append([]interface{}{ColorRed, ColorReset}, a...)...)
}
