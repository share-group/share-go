package server

type Server struct {
}

// 服务器接口
type IServer interface {
	// 设置接口前缀
	SetPrefix(prefix string)
	// 设置控制器入口
	SetController()
	// 启动http服务器
	Run()
}
