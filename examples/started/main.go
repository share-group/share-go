package main

import (
	_ "embed"

	"github.com/share-group/share-go/examples/started/controller/account"
	"github.com/share-group/share-go/examples/started/controller/system"
	"github.com/share-group/share-go/provider/formatter"
	"github.com/share-group/share-go/server"
)

//go:embed banner
var banner string

// 引入要暴露的控制器
var controllers = []any{
	account.AdminController,
	account.RoleController,
	system.CaptchaController,
}

func main() {
	s := server.NewHttpServer()
	s.SetBanner(banner)
	s.SetResponseFormatter(formatter.JSONResponseFormatter)
	s.RegisterControllers(controllers...)
	s.Run()
}
