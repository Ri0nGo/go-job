package log

import (
	"log/slog"
	"os"
)

func InitSlog(level slog.Level, source bool) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:     level,
		AddSource: source,
	}))

	slog.SetDefault(logger)
}
