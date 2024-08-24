package lsp

import (
	"bufio"
	"encoding/json"
	"io"
	"net/textproto"
	"strconv"
)

type Request struct {
	Version string           `json:"jsonrpc"`
	Id      *json.RawMessage `json:"id"`
	Method  string           `json:"method"`
	Params  json.RawMessage  `json:"params"`
}

func (r Request) IsJsonRPC() bool {
	return r.Version == "2.0"
}

func (r Request) IsNotification() bool {
	return r.Id == nil // lsp notifications never have an id
}

type ErrInvalidContentLengthHeader struct{}

func (e ErrInvalidContentLengthHeader) Error() string {
	return "Invalid content length header"
}

type ErrInvalidRequest struct{}

func (e ErrInvalidRequest) Error() string {
	return "Invalid lsp request"
}

func Read(r *bufio.Reader) (req Request, err error) {
	header, err := textproto.NewReader(r).ReadMIMEHeader()
	if err != nil {
		return
	}
	contentLength, err := strconv.ParseInt(header.Get("Content-Length"), 10, 64)
	if err != nil {
		return req, ErrInvalidContentLengthHeader{}
	}
	err = json.NewDecoder(io.LimitReader(r, contentLength)).Decode(&req)
	if err != nil {
		return
	}
	if !req.IsJsonRPC() {
		return req, ErrInvalidRequest{}
	}
	return
}
