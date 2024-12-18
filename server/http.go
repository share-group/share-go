package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/share-group/share-go/constant"
	HttpMethod "github.com/share-group/share-go/constant"
	"github.com/share-group/share-go/provider/config"
	"github.com/share-group/share-go/provider/formatter"
	loggerFactory "github.com/share-group/share-go/provider/logger"
	"github.com/share-group/share-go/provider/logging"
	"github.com/share-group/share-go/provider/validator"
	"github.com/share-group/share-go/util/maputil"
	"github.com/share-group/share-go/util/stringutil"
	"github.com/share-group/share-go/util/systemutil"
	"math"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var banner = ""

var handlers = make([]any, 0)

var urlMap = make(map[constant.HttpMethod][]string)

var middlewares = make([]func(next echo.HandlerFunc) echo.HandlerFunc, 0)

var responseFormatter func(fun func(c echo.Context) any) echo.HandlerFunc
var logger = loggerFactory.GetLogger()

type Server struct{}

// 设置打印banner
func (s *Server) SetBanner(bannerString string) {
	banner = bannerString
}

// 设置控制器入口
func (s *Server) RegisterControllers(controllers ...any) {
	handlers = append(handlers, controllers...)
}

// 设置中间件
func (s *Server) SetMiddlewares(middleware func(next echo.HandlerFunc) echo.HandlerFunc) {
	middlewares = append(middlewares, middleware)
}

// 设置返回数据格式器
func (s *Server) SetResponseFormatter(formatter func(fun func(c echo.Context) any) echo.HandlerFunc) {
	responseFormatter = formatter
}

// 启动服务器
func (s *Server) Run() {
	e := echo.New()
	addMiddleware(e)
	mappedHandler(e)
	showBanner()
	start(e)
}

func addMiddleware(e *echo.Echo) {
	for _, m := range middlewares {
		e.Use(m)
	}
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

	methods := maputil.Keys(methodFunMap)
	responseFormatter = systemutil.If(responseFormatter == nil, formatter.PlaintextResponseFormatter, responseFormatter).(func(fun func(c echo.Context) any) echo.HandlerFunc)

	// 自动注册路由
	for _, h := range handlers {
		obj := reflect.ValueOf(h)
		reflectType := reflect.TypeOf(h)
		for i := 0; i < reflectType.NumMethod(); i++ {
			m := reflectType.Method(i)

			// 如果注入了 echo.Context ，则把 echo.Context 反射进去：否则不需要
			paramIndex := math.MinInt32
			hasEchoContextIndex := false
			for j := 1; j < m.Type.NumIn(); j++ {
				tmpType := m.Type.In(j)
				if reflect.DeepEqual("echo.Context", tmpType.String()) {
					hasEchoContextIndex = true
				} else {
					paramIndex = j
					break
				}
			}

			// 约定路由规则，HttpMethod+接口名，例如：GetCaptcha，其实就是 GET /captcha；PostLogin，其实就是 POST /login，如果没有指定HttpMethod的话默认POST
			method := "POST"
			apiName := stringutil.FirstLowerCase(m.Name)
			prefix := config.GetString("server.prefix")
			module := strings.TrimSpace(strings.Split(fmt.Sprintf("%s", reflectType.Elem()), ".")[0])
			for _, httpMethod := range methods {
				httpMethod = stringutil.FirstUpperCase(strings.ToLower(httpMethod))
				if strings.HasPrefix(m.Name, httpMethod) {
					method = strings.ToUpper(httpMethod)
					apiName = stringutil.FirstLowerCase(m.Name[len(method):])
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
			logger.Info("%s %s %v", method, url, &m.Func)
			urlMap[HttpMethod.ValueOfHttpMethod(method)] = append(urlMap[HttpMethod.ValueOfHttpMethod(method)], url)

			var paramType reflect.Type
			if paramIndex != math.MinInt32 {
				paramType = m.Type.In(paramIndex)
			}

			// 注册路由方法
			methodFunMap[method](url, responseFormatter(func(c echo.Context) any {
				c.Set("requestTime", time.Now())
				callParam := []reflect.Value{obj}
				if hasEchoContextIndex {
					callParam = append(callParam, reflect.ValueOf(c))
				}
				body := validator.ValidateParameters(c, paramType)
				if body != nil {
					callParam = append(callParam, reflect.ValueOf(body))
				}
				logging.PrintRequestLog(c)

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
	logger.Info("%s server started on 0.0.0.0:%d in %s environment", config.GetString("application.name"), port, config.GetENV())
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
