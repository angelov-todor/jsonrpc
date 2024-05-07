package service

import (
	"fmt"
	"log"
)

type HelloService struct{}

type HelloRequest struct {
	Name string
	Age  int
}

type HelloResponse struct {
	Greeting string `json:"greeting"`
}

func (s *HelloService) Hello(req *HelloRequest, res *HelloResponse) error {
	log.Printf("Execute method: HelloService.Hello(); %v\n", req)
	res.Greeting = fmt.Sprintf("Hello: %s of Age: %d ! Welcome to the Go world!", req.Name, req.Age)

	return nil
}
