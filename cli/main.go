package main

import (
	"github.com/smtdfc/nagare/cli/cmd"
	"github.com/smtdfc/nagare/core"
)

func main() {
	core.SetupEnvironment()
	cmd.Execute()
}
