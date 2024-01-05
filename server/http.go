package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/share-group/share-go/bootstrap"
	"github.com/share-group/share-go/middleware"
	"github.com/share-group/share-go/util"
	"io"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var banner = ""
var handlers = make([]any, 0)
var logger = bootstrap.Logger.GetLogger()

type Server struct{}

// 设置打印banner
func (*Server) SetBanner(bannerString string) {
	banner = bannerString
}

// 设置处理器入口
func (*Server) SetHandlers(handler any) {
	handlers = append(handlers, handler)
}

func (*Server) Run() {
	e := echo.New()
	addMiddleware(e)
	mappedHandler(e)
	showBanner()
	start(e)
}

func addMiddleware(e *echo.Echo) {
	e.Use(middleware.Logger())
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

	// 是否启用数据验证器
	validatorEnable := bootstrap.Config.GetBoolValue("server.validator.enable")
	_validator := util.SystemUtil.If(validatorEnable, validator.New(), nil)

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
			prefix := bootstrap.Config.GetStringValue("server.prefix")
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
			methodFunMap[method](url, middleware.ResponseFormatter(func(c echo.Context) any {
				st := time.Now()
				callParam := []reflect.Value{obj}
				if paramType != nil {
					b, _ := io.ReadAll(c.Request().Body)
					body := reflect.New(paramType.Elem()).Interface()
					query := util.HttpUtil.ParseQueryString(c.Request().URL.String())
					json.Unmarshal(b, &body)
					json.Unmarshal([]byte(query), &body)
					bb, _ := json.Marshal(body)
					logger.Info(fmt.Sprintf("request %v, method: %v, data: %v", url, c.Request().Method, string(bb)))
					if validatorEnable {
						if err := _validator.(*validator.Validate).Struct(body); err != nil {
							panic(processErr(body, err))
							return nil
						}
					}
					callParam = append(callParam, reflect.ValueOf(body))
				}

				returnData := m.Func.Call(callParam)
				if len(returnData) == 1 || returnData[1].Interface() == nil {
					rr, _ := json.Marshal(returnData[0].Interface())
					logger.Info(fmt.Sprintf("url: %v, response: %v, data: %v, size: %v Byte, exec: %v", url, c.Response().Status, string(rr), len(rr), time.Since(st)))
					return returnData[0].Interface()
				}
				panic(returnData[1].Interface())
			}))
		}
	}
}

func processErr(obj interface{}, err error) string {
	if err == nil { //如果为nil 说明校验通过
		return ""
	}
	invalid, ok := err.(*validator.InvalidValidationError) //如果是输入参数无效，则直接返回输入参数错误
	if ok {
		return "输入参数错误：" + invalid.Error()
	}

	errorList := make([]string, 0)
	validationErrs := err.(validator.ValidationErrors) //断言是ValidationErrors
	for _, validationErr := range validationErrs {
		fieldName := validationErr.Field() //获取是哪个字段不符合格式
		typeOf := reflect.TypeOf(obj)
		// 如果是指针，获取其属性
		if typeOf.Kind() == reflect.Ptr {
			typeOf = typeOf.Elem()
		}
		field, _ := typeOf.FieldByName(fieldName) //通过反射获取filed
		message := strings.TrimSpace(fmt.Sprintf("%s", field.Tag.Get("message")))
		if len(message) > 0 {
			errorList = append(errorList, message)
		} else {
			errorList = append(errorList, strings.TrimSpace(fmt.Sprintf("%s", validationErr)))
		}
	}
	return strings.Join(errorList, ", ")
}

func showBanner() {
	if len(banner) > 0 {
		logger.Info(banner)
	}
}

func start(e *echo.Echo) {
	e.HidePort = true
	e.HideBanner = true
	port := bootstrap.Config.GetIntegerValue("server.port")
	if port <= 0 {
		logger.Fatal(fmt.Sprintf("invalid port: %d", port))
	}
	logger.Info(fmt.Sprintf("%s server started on 0.0.0.0:%d", bootstrap.Config.GetStringValue("application.name"), port))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
