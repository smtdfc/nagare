package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/smtdfc/nagare/core/utils"
)

var logger *slog.Logger

func init() {
	ts := time.Now().UnixNano()
	file, err := os.OpenFile(
		path.Join(utils.LogDir, fmt.Sprintf("nagare_%d.log", ts)),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		panic(err)
	}

	logger = slog.New(
		slog.NewJSONHandler(file, nil),
	)
}

func GetLogger() *slog.Logger {
	return logger
}
