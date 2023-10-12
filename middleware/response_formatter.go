package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"runtime/debug"
)

func ResponseFormatter(fun func(c echo.Context) any) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				fmt.Println(string(debug.Stack()))
				c.JSON(http.StatusInternalServerError, map[string]any{"code": 10001, "message": "服务器开了个小差"})
			}
		}()
		return c.JSON(http.StatusOK, map[string]any{"code": 0, "data": fun(c)})
	}
}
