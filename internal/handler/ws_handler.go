package handler

import (
	"bytes"
	"io"
	"log/slog"
	"net/rpc"
	"net/rpc/jsonrpc"

	"golang.org/x/net/websocket"
)

type WsHandler struct {
	rpc    *rpc.Server
	logger *slog.Logger
}

func NewWsHandler(logger *slog.Logger, rpc *rpc.Server) *WsHandler {
	return &WsHandler{
		rpc:    rpc,
		logger: logger,
	}
}

func (h *WsHandler) HandleRequest(ws *websocket.Conn) {
	for {
		var req []byte
		err := websocket.Message.Receive(ws, &req)
		if err != nil {
			h.logger.Error("ReadMessage:", "err", err)
			return
		}

		h.logger.Info("ServeRequest...")
		var res bytes.Buffer
		err = h.rpc.ServeRequest(jsonrpc.NewServerCodec(struct {
			io.ReadCloser
			io.Writer
		}{
			io.NopCloser(bytes.NewReader(req)),
			&res,
		}))
		if err != nil {
			h.logger.Error("ServeRequest:", "err", err)
			return
		}

		err = websocket.Message.Send(ws, res.Bytes())
		if err != nil {
			h.logger.Error("WriteMessage:", "err", err)
			return
		}
	}
}
