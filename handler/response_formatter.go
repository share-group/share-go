package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	exception "github.com/share-group/share-go/exception"
	loggerFactory "github.com/share-group/share-go/provider/logger"
	"net/http"
	"reflect"
	"runtime/debug"
)

var formatterLogger = loggerFactory.GetLogger("share.go.ResponseFormatter")

func ResponseFormatter(fun func(c echo.Context) any) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				formatterLogger.Info(fmt.Sprintf("\n%v\n%v\n", fmt.Sprint(err), string(debug.Stack())))
				formatterLogger.Info("xxxxxxxxxxxxxx  " + fmt.Sprintf("%v", c.Get("aaaaaaaaa")))
				if reflect.DeepEqual(reflect.TypeOf(err).String(), "errors.BusinessException") {
					e := err.(exception.BusinessException)
					c.JSON(http.StatusOK, map[string]any{"code": e.Code, "message": e.Message})
				} else {
					c.JSON(http.StatusInternalServerError, map[string]any{"code": 10001, "message": err})
				}
			}
		}()
		return c.JSON(http.StatusOK, map[string]any{"code": 0, "data": fun(c)})
	}
}
