package setup

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

var ServerVersion = "0.0.1"

func SetupLogging(logfileName string) *os.File {
	if !strings.Contains(logfileName, ".") {
		logfileName = logfileName + ".log"
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to get user home directory: %v", err)
	}
	logDir := filepath.Join(homeDir, ".local", "state")
	logFile := filepath.Join(logDir, "pygo-lsp.log")
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatalf("failed to create log directory: %v", err)
	}
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	return file
}
