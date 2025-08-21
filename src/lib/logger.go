package lib

import (
	"log/slog"
	"os"
)

func LoggerInit(debug bool) {

	var logLevel slog.Level

	if debug {
		logLevel = slog.LevelDebug
	} else {
		logLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   debug,
		Level:       logLevel,
		ReplaceAttr: nil,
	}))
	slog.SetDefault(logger)

	slog.Debug("Debug logging enabled.")
}
