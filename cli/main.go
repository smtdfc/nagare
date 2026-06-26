package main

import (
	"os"

	"github.com/smtdfc/nagare/cli/cmd"
	"github.com/smtdfc/nagare/core"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			core.GlobalPluginMgr.Shutdown()
			os.Exit(1)
		}
	}()
	core.SetupEnvironment()
	cmd.Execute()
}
