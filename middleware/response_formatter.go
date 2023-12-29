package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/share-group/share-go/bootstrap"
	"net/http"
	"runtime/debug"
)

func ResponseFormatter(fun func(c echo.Context) any) echo.HandlerFunc {
	logger := bootstrap.Logger.GetLogger()
	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				logger.Info(fmt.Sprint(err))
				logger.Info(string(debug.Stack()))
				c.JSON(http.StatusInternalServerError, map[string]any{"code": 10001, "message": err})
			}
		}()
		return c.JSON(http.StatusOK, map[string]any{"code": 0, "data": fun(c)})
	}
}
