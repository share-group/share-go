package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogURI:    true,
			LogStatus: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				zap.L().Info("request",
					zap.String("URI", v.URI),
					zap.Int("status", v.Status),
					zap.Any("headers", c.ParamValues()),
				)

				return nil
			},
		})

		//执行下一个中间件或者执行控制器函数, 然后返回执行结果
		return next(c)
	}
}
