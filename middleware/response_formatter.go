package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/share-group/share-go/bootstrap"
	exception "github.com/share-group/share-go/exception"
	"net/http"
	"reflect"
	"runtime/debug"
)

func ResponseFormatter(fun func(c echo.Context) any) echo.HandlerFunc {
	logger := bootstrap.Logger.GetLogger()
	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				logger.Info(fmt.Sprintf("\n%v\n%v\n", fmt.Sprint(err), string(debug.Stack())))
				if reflect.DeepEqual(reflect.TypeOf(err).String(), "errors.BusinessException") {
					e := err.(exception.BusinessException)
					c.JSON(http.StatusOK, map[string]any{"code": e.Code, "message": e.Message})
				} else {
					c.JSON(http.StatusInternalServerError, map[string]any{"code": 1, "message": err})
				}
			}
		}()
		return c.JSON(http.StatusOK, map[string]any{"code": 0, "data": fun(c)})
	}
}
