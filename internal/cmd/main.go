package main

import (
	"json-rpc/internal/handler"
	"json-rpc/internal/service"
	"log/slog"
	"net/http"
	"net/rpc"
	"os"

	"golang.org/x/net/websocket"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Starting http server")

	rpcServer := rpc.NewServer()

	rpcServer.Register(&service.HelloService{})
	rpcServer.Register(&service.TimeService{})

	wsHandler := handler.NewWsHandler(logger, rpcServer)
	httpHandler := handler.NewHttpHandler(logger, rpcServer)

	mux := http.NewServeMux()
	mux.Handle("/ws", websocket.Handler(wsHandler.HandleRequest))
	mux.Handle("/rpc", http.HandlerFunc(httpHandler.HandleRequest))

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("ListenAndServe:", "err", err)
		return
	}

	// http.ListenAndServe("localhost:8080", nil)
	// http.ListenAndServe(":8080", mux)
}
