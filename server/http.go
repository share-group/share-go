package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/share-group/share-go/provider/config"
	"github.com/share-group/share-go/provider/formatter"
	loggerFactory "github.com/share-group/share-go/provider/logger"
	"github.com/share-group/share-go/provider/logging"
	"github.com/share-group/share-go/provider/validator"
	"github.com/share-group/share-go/util"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var banner = ""

var handlers = make([]any, 0)

var responseFormatter func(fun func(c echo.Context) any) echo.HandlerFunc
var logger = loggerFactory.GetLogger()

type Server struct{}

// 设置打印banner
func (*Server) SetBanner(bannerString string) {
	banner = bannerString
}

// 设置处理器入口
func (*Server) SetHandlers(handler any) {
	handlers = append(handlers, handler)
}

// 设置返回数据格式器
func (*Server) SetResponseFormatter(formatter func(fun func(c echo.Context) any) echo.HandlerFunc) {
	responseFormatter = formatter
}

func (*Server) Run() {
	e := echo.New()
	addMiddleware(e)
	mappedHandler(e)
	showBanner()
	start(e)
}

func addMiddleware(e *echo.Echo) {
}

func mappedHandler(e *echo.Echo) {
	methodFunMap := map[string]func(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route{
		"GET":     e.GET,
		"HEAD":    e.HEAD,
		"POST":    e.POST,
		"PUT":     e.PUT,
		"DELETE":  e.DELETE,
		"CONNECT": e.CONNECT,
		"OPTIONS": e.OPTIONS,
		"TRACE":   e.TRACE,
		"PATCH":   e.PATCH,
	}

	responseFormatter = util.SystemUtil.If(responseFormatter == nil, formatter.PlaintextResponseFormatter, responseFormatter).(func(fun func(c echo.Context) any) echo.HandlerFunc)

	// 自动注册路由
	for _, h := range handlers {
		obj := reflect.ValueOf(h)
		reflectType := reflect.TypeOf(h)
		for i := 0; i < reflectType.NumMethod(); i++ {
			m := reflectType.Method(i)
			var paramType reflect.Type
			if m.Type.NumIn() > 1 {
				paramType = m.Type.In(1)
			}

			// 约定路由规则，HttpMethod+接口名，例如：GetCaptcha，其实就是 GET /captcha；PostLogin，其实就是 POST /login，如果没有指定HttpMethod的话默认POST
			method := "POST"
			apiName := util.StringUtil.FirstLowerCase(m.Name)
			prefix := config.GetString("server.prefix")
			module := strings.TrimSpace(strings.Split(fmt.Sprintf("%s", reflectType.Elem()), ".")[0])
			for _, httpMethod := range []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"} {
				httpMethod = util.StringUtil.FirstUpperCase(strings.ToLower(httpMethod))
				if strings.HasPrefix(m.Name, httpMethod) {
					method = strings.ToUpper(httpMethod)
					apiName = util.StringUtil.FirstLowerCase(m.Name[len(method):])
					break
				}
			}

			// 让url加多一层
			url := fmt.Sprintf("%s/%s/%s", prefix, module, apiName)
			controllerNameRegexp := regexp.MustCompile("\\(\\*(.*)\\.(.*)Controller")
			midUrl := strings.ReplaceAll(controllerNameRegexp.FindStringSubmatch(m.Func.String())[0], fmt.Sprintf("(*%s.", module), "")
			midUrl = strings.TrimSpace(strings.ReplaceAll(midUrl, "Controller", ""))
			if !strings.HasSuffix(url, midUrl) {
				url = fmt.Sprintf("%s/%s/%s/%s", prefix, module, midUrl, apiName)
			}
			logger.Info(fmt.Sprintf("%s %s %v", method, strings.ReplaceAll(url, prefix, ""), &m.Func))

			// 注册路由方法
			methodFunMap[method](url, responseFormatter(func(c echo.Context) any {
				c.Set("requestTime", time.Now())
				callParam := []reflect.Value{obj}
				body := validator.ValidateParameters(c, paramType)
				if body != nil {
					callParam = append(callParam, reflect.ValueOf(body))
				}
				go logging.PrintRequestLog(c)

				// 约定，只有一个返回，或者没有
				returnData := m.Func.Call(callParam)
				if len(returnData) <= 0 {
					return nil
				}

				return returnData[0].Interface()
			}))
		}
	}
}

func showBanner() {
	if len(banner) > 0 {
		logger.Info(banner)
	}
}

func start(e *echo.Echo) {
	e.HidePort = true
	e.HideBanner = true
	port := config.GetInt("server.port")
	if port <= 0 {
		logger.Fatal(fmt.Sprintf("invalid port: %d", port))
	}
	logger.Info(fmt.Sprintf("%s server started on 0.0.0.0:%d", config.GetString("application.name"), port))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
