package mux

type TextDocumentSyncKind int

const (
	TextDocumentSyncKindNone TextDocumentSyncKind = iota
	TextDocumentSyncKindFull
	TextDocumentSyncKindIncremental
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
