package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func NewHttpServer() *Server {
	return &Server{}
}

func (*Server) Run() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
