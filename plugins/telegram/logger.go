package main

import (
	"log/slog"

	nagare_logger "github.com/smtdfc/nagare/shared/logger"
	nagare_path "github.com/smtdfc/nagare/shared/path"
)

var PluginLogger *slog.Logger

func init() {
	nagare_logger.InitLogger(nagare_path.GetPluginLogDir("org.smtdfc.nagare.plugins.telegram"))
	PluginLogger = nagare_logger.GetLogger("org.smtdfc.nagare.plugins.telegram")
}
