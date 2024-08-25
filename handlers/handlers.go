package handlers

import (
	"github.com/jesses-code-adventures/pygo-lsp/mux"
)

func RegisterHandlers(m *mux.Mux) {
	initialize(m)
	initialized(m)
}
