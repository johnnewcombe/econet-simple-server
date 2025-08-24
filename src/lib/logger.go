package lib

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
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

func LogDebugData(data []byte) {

	if slog.Default().Enabled(context.TODO(), slog.LevelDebug) && len(data) > 0 {
		dump := strings.Split(HexDump(data), "\n")
		for i, line := range dump {
			if i == 0 {
				fmt.Println("\tdata=" + line)
			} else {
				fmt.Println("\t     " + line)
			}
		}
	}
}
