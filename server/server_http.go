package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	logger "github.com/share-group/share-go/bootstrap"
	"net/http"
)

func NewHttpServer() *Server {
	return &Server{}
}

// 设置接口前缀
func (*Server) SetPrefix(prefix string) {}

// 设置控制器入口
func (*Server) SetController() {}

func (*Server) Run() {
	e := echo.New()
	//if l, ok := e.Logger.(*logger.Logger); ok {
	//	logFile, err := os.OpenFile(path.Join(cmd, "log/log.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	//	if err != nil {
	//		fmt.Println("open log file failed, err:", err)
	//		return
	//	}
	//	l.SetHeader("${time_rfc3339} ${level}")
	//	l.SetOutput(logFile)
	//}
	e.Use(middleware.Logger())
	e.GET("/", func(c echo.Context) error {
		logger.Logger.Info("哈哈哈哈哈")
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
