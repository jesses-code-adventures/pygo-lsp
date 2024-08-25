package mux

import (
	"github.com/jesses-code-adventures/pygo-lsp/setup"
	"log/slog"
	"os"
	"testing"
)

func TestMuxInitialization(t *testing.T) {
	logfile := setup.SetupLogging("pygo-lsp-test.log")
	mux := NewMux(os.Stdin, os.Stdout, setup.ServerVersion, slog.New(slog.NewJSONHandler(logfile, nil)))
	if len(mux.notificationHandlers) == 0 {
		t.Errorf("Expected notificationHandlers length to be > 0")
	}
	if len(mux.methodHandlers) == 0 {
		t.Errorf("Expected methodHandlers length to be > 0")
	}
}
