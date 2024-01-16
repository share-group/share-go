package main

import (
	"github.com/share-group/share-go/examples/started/controller/account"
	"github.com/share-group/share-go/examples/started/controller/system"
	"github.com/share-group/share-go/provider/formatter"
	"github.com/share-group/share-go/server"
)

func main() {
	s := server.NewHttpServer()
	s.SetResponseFormatter(formatter.JSONResponseFormatter)
	s.SetHandlers(account.AdminController)
	s.SetHandlers(system.CaptchaController)
	s.Run()
}
