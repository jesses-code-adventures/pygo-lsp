package handlers

import (
	"encoding/json"

	"github.com/jesses-code-adventures/pygo-lsp/mux"
)

type ClientCapabilities struct{}

type TextDocumentSyncKind int

const (
	TextDocumentSyncKindNone TextDocumentSyncKind = iota
	TextDocumentSyncKindFull
	TextDocumentSyncKindIncremental
)

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

type InitializeParams struct {
	ClientInfo   *ClientInfo        `json:"clientInfo"`
	Capabilities ClientCapabilities `json:"capabilities"`
}

type ClientInfo struct {
	Name    string  `json:"name"`
	Version *string `json:"version"`
}

type ServerCapabilities struct {
	TextDocumentSync TextDocumentSyncKind `json:"textDocumentSync"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func initialize(m *mux.Mux) {
	m.HandleMethod("initialize", func(params json.RawMessage) (result any, err error) {
		var initializeParams InitializeParams
		if err = json.Unmarshal(params, &initializeParams); err != nil {
			return result, err
		}
		result = InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync: TextDocumentSyncKindFull,
			},
			ServerInfo: ServerInfo{Name: "pygo-lsp", Version: m.ServerVersion},
		}
		return result, err
	})
}
