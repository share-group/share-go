package server

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	config "github.com/share-group/share-go/bootstrap"
	"github.com/share-group/share-go/middleware"
	"go.uber.org/zap"
	"io"
	"reflect"
	"strings"
)

var handlers = make([]any, 0)

type Server struct {
}

// 设置处理器入口
func (*Server) SetHandlers(handler any) {
	handlers = append(handlers, handler)
}

func (*Server) Run() {
	e := echo.New()
	e.Use(middleware.Logger)
	for _, handler := range handlers {
		obj := reflect.ValueOf(handler)
		reflectType := reflect.TypeOf(handler)
		for i := 0; i < reflectType.NumMethod(); i++ {
			m := reflectType.Method(i)
			paramType := m.Type.In(1)
			method := "POST"
			module := strings.TrimSpace(strings.Split(fmt.Sprintf("%s", reflectType.Elem()), ".")[0])
			url := fmt.Sprintf("/%s/%s", module, strings.ToLower(m.Name))
			zap.L().Info(fmt.Sprintf("register url: [%s] %s", method, url))
			e.POST(url, middleware.ResponseFormatter(func(c echo.Context) any {
				b, _ := io.ReadAll(c.Request().Body)
				body := reflect.New(paramType).Interface()
				json.Unmarshal(b, &body)
				return m.Func.Call([]reflect.Value{obj, reflect.ValueOf(body).Elem()})[0].Interface()
			}))
		}
	}

	start(e)
}

func start(e *echo.Echo) {
	e.HidePort = true
	e.HideBanner = true
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", config.Get("server.port"))))
}
