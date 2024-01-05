package main

import (
	"github.com/share-group/share-go/examples/started/controller/account"
	"github.com/share-group/share-go/examples/started/controller/system"
	"github.com/share-group/share-go/server"
)

func main() {
	s := server.NewHttpServer()
	s.SetHandlers(account.AdminController)
	s.SetHandlers(system.CaptchaController)
	s.Run()
}
