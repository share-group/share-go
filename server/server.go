package server

func NewHttpServer() *Server {
	return &Server{}
}

// 服务器接口
type IServer interface {
	// 设置打印banner
	SetBanner(banner string)
	// 设置处理器入口
	SetHandlers(handler any)
	// 启动http服务器
	Run()
}
