package lsp

import (
	"bufio"
	"bytes"
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

type LspHeaders struct {
	ContentLength int64
}

func ReadHeaders(r *bufio.Reader) (headers LspHeaders, err error) {
	header, err := textproto.NewReader(r).ReadMIMEHeader()
	if err != nil {
		return
	}
	contentLength, err := strconv.ParseInt(header.Get("Content-Length"), 10, 64)
	if err != nil {
		return headers, ErrInvalidContentLengthHeader{}
	}
	return LspHeaders{ContentLength: contentLength}, err
}

func Read(r *bufio.Reader, headers LspHeaders) (req Request, err error) {
	err = json.NewDecoder(io.LimitReader(r, headers.ContentLength)).Decode(&req)
	if err != nil {
		return
	}
	if !req.IsJsonRPC() {
		return req, ErrInvalidRequest{}
	}
	return
}

func RawBytes(r *bufio.Reader, headers LspHeaders) (b []byte, err error) {
	var buf bytes.Buffer
	tee := io.TeeReader(io.LimitReader(r, headers.ContentLength), &buf)
	tempBuf := make([]byte, headers.ContentLength)
	_, err = io.ReadFull(tee, tempBuf)
	if err != nil {
		return b, err
	}
	b = make([]byte, headers.ContentLength)
	copy(b, tempBuf)
	r.Reset(io.MultiReader(&buf, r))
	return b, err
}
