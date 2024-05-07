package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/rpc"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type HttpHandler struct {
	rpc    *rpc.Server
	logger *slog.Logger
}

func NewHttpHandler(logger *slog.Logger, rpc *rpc.Server) *HttpHandler {
	return &HttpHandler{
		rpc:    rpc,
		logger: logger,
	}
}

func (h *HttpHandler) HandleRequest(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ct := req.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	// r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB limit

	h.logger.Info("ServeRequest...")

	codec, err := newHttpCodec(req, w)
	if err != nil {
		h.logger.Error("NewHttpCodec:", "err", err)
		return
	}
	err = h.rpc.ServeRequest(codec)
	if err != nil {
		h.logger.Error("ServeRequest:", "err", err)
		return
	}
}

type httpCodec struct {
	writer     http.ResponseWriter
	bodyBuf    *bytes.Buffer
	rpcRequest *RPCRequest
}

func newHttpCodec(r *http.Request, w http.ResponseWriter) (*httpCodec, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode request header")
	}
	defer r.Body.Close()
	bodyBuf := bytes.NewBuffer(body)
	if bodyBuf == nil {
		return nil, errors.New("empty request body")
	}

	rpcReq := &RPCRequest{}
	dec := json.NewDecoder(bodyBuf)
	err = dec.Decode(rpcReq)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode request")
	}

	return &httpCodec{
		writer:     w,
		bodyBuf:    bodyBuf,
		rpcRequest: rpcReq,
	}, nil
}

func (c *httpCodec) ReadRequestHeader(r *rpc.Request) error {
	r.ServiceMethod = c.rpcRequest.Method
	// r.Seq = c.rpcRequest.ID
	return nil
}

func (c *httpCodec) ReadRequestBody(x interface{}) error {
	buf := bytes.NewBuffer([]byte{})
	// encode the request params as bytes
	err := json.NewEncoder(buf).Encode(c.rpcRequest.Params)
	if err != nil {
		return errors.Wrap(err, "could not encode request body params")
	}

	err = json.NewDecoder(buf).Decode(x)
	if err != nil {
		return errors.Wrap(err, "could not decode request body")
	}
	return nil
}

func (c *httpCodec) WriteResponse(r *rpc.Response, x interface{}) error {
	c.writer.Header().Set("Content-Type", "application/json")

	reply := &RPCResponse{
		JSONRPC: "2.0",
		Result:  x,
		ID:      strconv.FormatUint(r.Seq, 10),
		Error:   nil,
	}
	err := json.NewEncoder(c.writer).Encode(reply)
	if err != nil {
		return errors.Wrap(err, "could not encode response")
	}
	return nil
}

func (c *httpCodec) Close() error {
	return nil
}
