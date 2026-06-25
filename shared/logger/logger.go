package nagare_logger

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"time"
)

var Logger *slog.Logger

func InitLogger(logDir string) error {
	ts := time.Now().UnixNano()
	file, err := os.OpenFile(
		path.Join(logDir, fmt.Sprintf("nagare_%d.log", ts)),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		return fmt.Errorf("could not open log file: %w", err)
	}

	Logger = slog.New(
		slog.NewJSONHandler(file, nil),
	)

	return nil
}

func GetLogger(name string) *slog.Logger {
	return Logger.With(slog.String("module", name))
}
