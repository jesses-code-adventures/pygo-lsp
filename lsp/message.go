package lsp

type Message interface {
	IsJsonRPC() bool
}

type Notification struct {
	ProtocolVersion string `json:"jsonrpc"`
	Method          string `json:"method"`
	Params          any    `json:"params"`
}

func (n Notification) IsJsonRPC() bool {
	return n.ProtocolVersion == "2.0"
}
