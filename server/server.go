package server

type Server struct {
}

// 服务器接口
type IServer interface {
	// 启动http服务器
	Run()
}
