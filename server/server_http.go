package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"path"
)

func NewHttpServer() *Server {
	return &Server{}
}

func (*Server) Run() {
	logger.NewLogger()
	e := echo.New()
	if l, ok := e.Logger.(*logger.Logger); ok {
		logFile, err := os.OpenFile(path.Join(cmd, "log/log.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println("open log file failed, err:", err)
			return
		}
		l.SetHeader("${time_rfc3339} ${level}")
		l.SetOutput(logFile)
	}
	e.Use(middleware.Logger())
	e.GET("/", func(c echo.Context) error {
		logger.Info("哈哈哈哈哈")
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
