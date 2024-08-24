package lsp

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"
)

type NotificationHandler func(params json.RawMessage) (err error)
type MethodHandler func(params json.RawMessage) (result any, err error)

type Mux struct {
	reader               *bufio.Reader
	writer               *bufio.Writer
	notificationHandlers map[string]NotificationHandler
	methodHandlers       map[string]MethodHandler
	writeLock            *sync.Mutex
}

func NewMux(r *os.File, w *os.File) Mux {
	mut := sync.Mutex{}
	return Mux{
		reader:               bufio.NewReader(r),
		writer:               bufio.NewWriter(w),
		writeLock:            &mut,
		notificationHandlers: make(map[string]NotificationHandler),
		methodHandlers:       make(map[string]MethodHandler),
	}
}

type ErrMethodNotFound struct{}

func (e ErrMethodNotFound) Error() string { return "Method not found" }

func (m *Mux) processSingle(req Request) (err error) {
	if req.IsNotification() {
		if nh, ok := m.notificationHandlers[req.Method]; ok {
			return nh(req.Params)
		}
		return
	}
	mh, ok := m.methodHandlers[req.Method]
	if !ok {
		return m.write(NewErrorResponse(req.Id, ErrMethodNotFound{}))
	}
	result, err := mh(req.Params)
	if err != nil {
		return m.write(NewErrorResponse(req.Id, err))
	}
	return m.write(NewResponse(req.Id, result))
}

func (m *Mux) Process() (err error) {
	req, err := Read(m.reader)
	if err != nil {
		return err
	}
	go func(req Request) {
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

func (m *Mux) write(msg Message) (err error) {
	m.writeLock.Lock()
	defer m.writeLock.Unlock()
	return Write(m.writer, msg)
}

func (m *Mux) Notify(method string, params any) (err error) {
	n := Notification{
		ProtocolVersion: "2.0",
		Method:          method,
		Params:          params,
	}
	return m.write(n)
}
