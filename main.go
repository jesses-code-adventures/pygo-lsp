package main

import (
	"log/slog"
	"os"

	"github.com/jesses-code-adventures/pygo-lsp/mux"
	"github.com/jesses-code-adventures/pygo-lsp/setup"
)

func main() {
	logfile := setup.SetupLogging("pygo-lsp.log")
	defer logfile.Close()
	logger := slog.New(slog.NewJSONHandler(logfile, nil))
	m := mux.NewMux(os.Stdin, os.Stdout, setup.ServerVersion, logger)
	m.Run()
}
