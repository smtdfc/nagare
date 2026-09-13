package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/smtdfc/nagare/shared/paths"
	"gopkg.in/natefinch/lumberjack.v2"
)

type BaseLogger struct {
	slog.Logger
}

func (b *BaseLogger) Clone() *BaseLogger {
	cloned := b.Logger.With()
	return &BaseLogger{
		Logger: *cloned,
	}
}

func (b *BaseLogger) With(a ...any) *BaseLogger {
	return &BaseLogger{
		Logger: *b.Logger.With(a...),
	}
}

// @Injectable
func NewBaseLogger() (*BaseLogger, error) {
	if err := os.MkdirAll(paths.LogDir, 0755); err != nil {
		return nil, err
	}

	logPath := filepath.Join(paths.LogDir, "app.log")

	fileRotator := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100,
		MaxBackups: 30,
		MaxAge:     30,
		Compress:   true,
	}

	multiWriter := io.MultiWriter(os.Stdout, fileRotator)
	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	baseLogger := slog.New(handler)

	return &BaseLogger{
		Logger: *baseLogger,
	}, nil
}
