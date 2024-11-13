package server

import (
	"github.com/labstack/echo/v4"
)

func NewHttpServer() *Server {
	return &Server{}
}

// 服务器接口
type IServer interface {
	// 设置打印banner
	SetBanner(banner string)
	// 设置控制器入口
	RegisterControllers(controllers ...any)
	// 设置中间件
	SetMiddlewares(middleware func(next echo.HandlerFunc) echo.HandlerFunc)
	// 设置返回数据格式器
	SetResponseFormatter(func(fun func(c echo.Context) any) echo.HandlerFunc)
	// 启动http服务器
	Run()
}
