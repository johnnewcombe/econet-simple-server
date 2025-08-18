package main

import (
	"log/slog"
	"os"

	"github.com/johnnewcombe/econet-simple-server/src/cobra"
)

var Logger *slog.Logger

func main() {
	cobra.Execute()
	os.Exit(0)
}
