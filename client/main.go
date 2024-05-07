package main

import (
	"bufio"
	"json-rpc/internal/service"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"

	"golang.org/x/net/websocket"
)

func call(c *rpc.Client, name string) {
	req := service.HelloRequest{Name: name, Age: 30}
	var res service.HelloResponse

	err := c.Call("HelloService.Hello", req, &res)
	if err != nil {
		log.Fatal("error:", err)
	}
	log.Printf("Response: %s", res.Greeting)
}

func main() {
	ws, err := websocket.Dial("ws://localhost:8080/ws", "", "http://localhost/")
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	c := jsonrpc.NewClient(ws)

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		call(c, text)
	}
}
