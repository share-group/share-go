package server

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/share-group/share-go/provider/config"
)

type DaemonServer struct{}

func NewDaemonServer() *DaemonServer {
	return &DaemonServer{}
}

// 设置打印banner
func (s *DaemonServer) SetBanner(bannerString string) {
	banner = bannerString
}

// 设置控制器入口
func (s *DaemonServer) RegisterControllers(controllers ...any) {
	logger.Fatal("not supported RegisterControllers")
}

// 设置中间件
func (s *DaemonServer) SetMiddlewares(middleware func(next echo.HandlerFunc) echo.HandlerFunc) {
	logger.Fatal("not supported SetMiddlewares")
}

// 设置返回数据格式器
func (s *DaemonServer) SetResponseFormatter(formatter func(fun func(c echo.Context) any) echo.HandlerFunc) {
	logger.Fatal("not supported SetResponseFormatter")
}

// 启动服务器
func (s *DaemonServer) Run() {
	showBanner()
	logger.Info("%s server started in %s environment", config.GetString("application.name"), config.GetENV())
	<-context.Background().Done() // 阻塞程序，直到收到系统信号
}
