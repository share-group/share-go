package main

import (
	"github.com/share-group/share-go/server"
)

func main() {
	s := server.NewHttpServer()
	s.SetPrefix("/api/v1")
	s.SetController()
	s.Run()
}
