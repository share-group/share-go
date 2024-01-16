package server

import "github.com/labstack/echo/v4"

func NewHttpServer() *Server {
	return &Server{}
}

// 服务器接口
type IServer interface {
	// 设置打印banner
	SetBanner(banner string)
	// 设置处理器入口
	SetHandlers(handler any)
	// 设置返回数据格式器
	SetResponseFormatter(func(fun func(c echo.Context) any) echo.HandlerFunc)
	// 启动http服务器
	Run()
}
