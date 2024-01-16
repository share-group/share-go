package formatter

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	exception "github.com/share-group/share-go/exception"
	loggerFactory "github.com/share-group/share-go/provider/logger"
	"github.com/share-group/share-go/provider/logging"
	"net/http"
	"reflect"
	"runtime/debug"
)

var logger = loggerFactory.GetLogger()

// 明文数据格式
func PlaintextResponseFormatter(fun func(c echo.Context) any) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				var response any
				logger.Info(fmt.Sprintf("\n%v\n%v\n", fmt.Sprint(err), string(debug.Stack())))
				if reflect.DeepEqual(reflect.TypeOf(err).String(), "errors.BusinessException") {
					e := err.(exception.BusinessException)
					response = e.Message
					c.JSON(http.StatusOK, response)
				} else {
					response = err
					c.JSON(http.StatusInternalServerError, response)
				}
				c.Set("response", []byte(fmt.Sprintf("%v", response)))
				go logging.SaveStringRequestLog(c)
			}
		}()

		response := fun(c)
		b, _ := json.Marshal(fun(c))
		c.Set("response", b)
		go logging.SaveStringRequestLog(c)
		return c.String(http.StatusOK, fmt.Sprintf("%v", response))
	}
}

// json结构的数据格式
func JSONResponseFormatter(fun func(c echo.Context) any) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				var response map[string]any
				logger.Info(fmt.Sprintf("\n%v\n%v\n", fmt.Sprint(err), string(debug.Stack())))
				if reflect.DeepEqual(reflect.TypeOf(err).String(), "errors.BusinessException") {
					e := err.(exception.BusinessException)
					response = map[string]any{"code": e.Code, "message": e.Message}
					c.JSON(http.StatusOK, response)
				} else {
					response = map[string]any{"code": 10001, "message": err}
					c.JSON(http.StatusInternalServerError, response)
				}
				b, _ := json.Marshal(response)
				c.Set("response", b)
				go logging.SaveJSONRequestLog(c)
			}
		}()

		response := fun(c)
		b, _ := json.Marshal(fun(c))
		c.Set("response", b)
		go logging.SaveJSONRequestLog(c)
		return c.JSON(http.StatusOK, map[string]any{"code": 0, "data": response})
	}
}
