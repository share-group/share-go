package main

import (
	system "github.com/share-group/share-go/examples/started/system/controller"
	"github.com/share-group/share-go/server"
)

func main() {
	s := server.NewHttpServer()
	s.SetHandlers(system.SystemController)
	s.Run()
}
