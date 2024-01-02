package middleware

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/share-group/share-go/bootstrap"
)

func Logger() echo.MiddlewareFunc {
	logger := bootstrap.Logger.GetLogger()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()
			logger.Info(fmt.Sprintf("%v %v", request.Method, request.URL))
			return nil
		}
	}
}
