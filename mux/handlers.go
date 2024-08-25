package mux

import (
	"encoding/json"
	"fmt"
	"time"
)

func (m *Mux) RegisterHandlers() {
	m.HandleMethod("initialize", func(params json.RawMessage) (result any, err error) {
		var initializeParams InitializeParams
		if err = json.Unmarshal(params, &initializeParams); err != nil {
			return result, err
		}
		result = InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync: TextDocumentSyncKindFull,
			},
			ServerInfo: ServerInfo{Name: "pygo-lsp", Version: m.serverVersion},
		}
		return result, err
	})
	m.HandleNotification("initialized", func(params json.RawMessage) (err error) {
		go func() {
			count := 1
			for {
				time.Sleep(time.Second * 1)
				m.Notify("window/showMessage", ShowMessageParams{
					Type:    MessageTypeInfo,
					Message: fmt.Sprintf("shown %d messages", count),
				})
				count++
			}
		}()
		return err
	})
}
