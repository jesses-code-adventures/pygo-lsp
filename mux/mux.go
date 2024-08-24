package mux

import (
	"bufio"
	"encoding/json"
	"log/slog"
	"os"
	"sync"

	"github.com/jesses-code-adventures/pygo-lsp/lsp"
)

type NotificationHandler func(params json.RawMessage) (err error)
type MethodHandler func(params json.RawMessage) (result any, err error)

type Mux struct {
	reader               *bufio.Reader
	writer               *bufio.Writer
	notificationHandlers map[string]NotificationHandler
	methodHandlers       map[string]MethodHandler
	writeLock            *sync.Mutex
	serverVersion        string
	logger               *slog.Logger
}

func NewMux(r *os.File, w *os.File, version string, logger *slog.Logger) Mux {
	mut := sync.Mutex{}
	return Mux{
		reader:               bufio.NewReader(r),
		writer:               bufio.NewWriter(w),
		writeLock:            &mut,
		notificationHandlers: make(map[string]NotificationHandler),
		methodHandlers:       make(map[string]MethodHandler),
		serverVersion:        version,
		logger:               logger,
	}
}

type ErrMethodNotFound struct{}

func (e ErrMethodNotFound) Error() string { return "Method not found" }

func (m *Mux) processSingle(req lsp.Request) (err error) {
	if req.IsNotification() {
		if nh, ok := m.notificationHandlers[req.Method]; ok {
			return nh(req.Params)
		}
		return
	}
	mh, ok := m.methodHandlers[req.Method]
	if !ok {
		return m.write(lsp.NewErrorResponse(req.Id, ErrMethodNotFound{}))
	}
	result, err := mh(req.Params)
	if err != nil {
		return m.write(lsp.NewErrorResponse(req.Id, err))
	}
	return m.write(lsp.NewResponse(req.Id, result))
}

func (m *Mux) Process() (err error) {
	req, err := lsp.Read(m.reader)
	if err != nil {
		return err
	}
	go func(req lsp.Request) {
		m.processSingle(req)
	}(req)
	return
}

func (m *Mux) HandleMethod(name string, method MethodHandler) {
	m.methodHandlers[name] = method
}

func (m *Mux) HandleNotification(name string, notification NotificationHandler) {
	m.notificationHandlers[name] = notification
}

func (m *Mux) write(msg lsp.Message) (err error) {
	m.writeLock.Lock()
	defer m.writeLock.Unlock()
	return lsp.Write(m.writer, msg)
}

func (m *Mux) Notify(method string, params any) (err error) {
	n := lsp.Notification{
		ProtocolVersion: "2.0",
		Method:          method,
		Params:          params,
	}
	return m.write(n)
}

func (m *Mux) Run() {
	for {
		if err := m.Process(); err != nil {
			return
		}
	}
}
