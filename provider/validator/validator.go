package validator

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	exception "github.com/share-group/share-go/exception"
	"github.com/share-group/share-go/provider/config"
	"github.com/share-group/share-go/util/httputil"
	"github.com/share-group/share-go/util/systemutil"
	"io"
	"reflect"
	"strings"
)

var enable = config.GetBool("server.validator.enable")
var _validator = systemutil.If(enable, validator.New(), nil)

func ValidateParameters(c echo.Context, paramType reflect.Type) any {
	if paramType == nil {
		return nil
	}

	b, _ := io.ReadAll(c.Request().Body)
	body := reflect.New(paramType.Elem()).Interface()
	query := httputil.ParseQueryString(c.Request().URL.String())
	json.Unmarshal(b, &body)
	json.Unmarshal([]byte(query), &body)
	request, _ := json.Marshal(body)
	c.Set("request", request)
	if !enable {
		return body
	}

	if err := _validator.(*validator.Validate).Struct(body); err != nil {
		panic(exception.NewBusinessException(10002, processErr(body, err)))
		return nil
	}
	return body
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
