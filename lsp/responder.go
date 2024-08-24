package lsp

import (
	"bufio"
	"encoding/json"
	"fmt"
)

type Response struct {
	ProtocolVersion string           `json:"jsonrpc"`
	Id              *json.RawMessage `json:"id"`
	Result          any              `json:"result"`
	Error           *Error           `json:"error"`
}

func NewResponse(id *json.RawMessage, result any) Response {
	return Response{
		ProtocolVersion: "2.0",
		Id:              id,
		Result:          result,
		Error:           nil,
	}
}

func NewErrorResponse(id *json.RawMessage, err error) Response {
	encodedError := NewError(err)
	return Response{
		ProtocolVersion: "2.0",
		Id:              id,
		Result:          nil,
		Error:           &encodedError,
	}
}

func (r Response) IsJsonRPC() bool { return r.ProtocolVersion == "2.0" }

type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewError(err error) Error {
	return Error{
		Code:    0,
		Message: err.Error(),
		Data:    nil,
	}
}

func (e *Error) Error() string { return e.Message }

func Write(w *bufio.Writer, msg Message) (err error) {
	body, err := json.Marshal(msg)
	if err != nil {
		return
	}
	headers := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(body))
	if _, err = w.WriteString(headers); err != nil {
		return
	}
	if _, err = w.Write(body); err != nil {
		return
	}
	return w.Flush()
}
