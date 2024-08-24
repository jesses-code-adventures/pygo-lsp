package main

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/jesses-code-adventures/pygo-lsp/mux"
)

func main() {
	serverVersion := "0.0.1"
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	logDir := filepath.Join(homeDir, ".local", "state")
	logFile := filepath.Join(logDir, "pygo-lsp.log")
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		panic(err)
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	logger := slog.New(slog.NewJSONHandler(file, nil))
	m := mux.NewMux(os.Stdin, os.Stdout, serverVersion, logger)
	m.Run()
}
