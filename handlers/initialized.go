package handlers

import (
	"encoding/json"

	"github.com/jesses-code-adventures/pygo-lsp/mux"
)

type MessageType int

const (
	NoMessageType MessageType = iota
	MessageTypeError
	MessageTypeWarining
	MessageTypeInfo
	MessageTypeLog
	MessageTypeDebug
)

type ShowMessageParams struct {
	Type    MessageType `json:"type"`
	Message string      `json:"message"`
}

func initialized(m *mux.Mux) {
	m.HandleNotification("initialized", func(params json.RawMessage) (err error) {
		go func() {
			m.Notify("window/showMessage", ShowMessageParams{
				Type:    MessageTypeInfo,
				Message: "welcome to pygo-lsp :-)",
			})
		}()
		return err
	})
}
