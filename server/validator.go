package server

//
//import (
//	"910.com/plus2.git/provider/invocation"
//	"910.com/plus2.git/routing"
//	"910.com/plus2.git/utils/map_util"
//	"910.com/plus2.git/utils/string_util"
//	"errors"
//	"fmt"
//	vld "github.com/go-playground/validator/v10"
//	"log"
//	"reflect"
//	"strconv"
//	"strings"
//)
//
//type ValidatorExclude struct {
//	objectName string
//	mathodName string
//}
//
//type Validator struct {
//	v        *vld.Validate
//	enable   bool
//	excludes []*ValidatorExclude
//}
//
///**
// * 构造函数
// * @param Validation 自定义校验规则
// */
//func NewValidator(validation any) *Validator {
//	v := vld.New()
//	if validation == nil {
//		return &Validator{v: v, excludes: make([]*ValidatorExclude, 0)}
//	}
//	obj := reflect.ValueOf(validation)
//	reflectType := reflect.TypeOf(validation)
//	for i := 0; i < reflectType.NumMethod(); i++ {
//		method := reflectType.Method(i)
//		tag := string_util.LcFirst(method.Name)
//		v.RegisterValidation(tag, func(f vld.FieldLevel) bool {
//			return method.Func.Call([]reflect.Value{obj, reflect.ValueOf(f)})[0].Bool()
//		})
//		log.Printf("成功加载自定义校验规则：%s", tag)
//	}
//	return &Validator{v: v, excludes: make([]*ValidatorExclude, 0)}
//}
//
///**
// * 注入验证器配置
// * @param config 配置信息
// */
//func (v *Validator) SetConfig(config map[string]interface{}) error {
//	// 解析配置文件
//	v.enable, _ = strconv.ParseBool(fmt.Sprintf("%v", map_util.GetValueFromMap(config, "enable", "false")))
//	conf := strings.Split(strings.TrimSpace(fmt.Sprintf("%v", map_util.GetValueFromMap(config, "excludes", ""))), ",")
//	if len(conf) <= 1 && len(conf[0]) <= 0 {
//		return nil
//	}
//
//	// 为了以防开发人员手误，输入：goods.List,user.*,goods.* 这种规则，所以需要做一下处理，同名的objectName和methodName规则只取第一个
//	objectNameMap := make(map[string]interface{})
//	mathodNameMap := make(map[string]interface{})
//	for _, c := range conf {
//		// 如果不是按照 objectName.methodName 的格式，则忽略
//		arr := strings.Split(strings.TrimSpace(c), ".")
//		if len(arr) != 2 {
//			log.Println(fmt.Sprintf("validator.default.excludes 数据校验器白名单规则[%s]格式错误，此配置项不生效。正确格式为：objectName.methodName", c))
//			continue
//		}
//
//		objectName := strings.TrimSpace(arr[0])
//		mathodName := strings.TrimSpace(arr[1])
//		if map_util.ContainsKey(objectNameMap, objectName) && map_util.ContainsKey(mathodNameMap, mathodName) {
//			log.Println(fmt.Sprintf("validator.default.excludes 数据校验器白名单规则[%s]已存在，此配置项不生效", c))
//			continue
//		}
//		objectNameMap[objectName] = true
//		mathodNameMap[mathodName] = true
//		v.excludes = append(v.excludes, &ValidatorExclude{objectName: objectName, mathodName: mathodName})
//	}
//
//	return nil
//}
//
///**
// * 通过请求上下文验证数据
// * @param ctx 请求上下文
// * @param obj 需要验证的数据
// */
//func (v *Validator) ValidatCtx(ctx *routing.Context, obj interface{}) error {
//	if v.isExclude(ctx) {
//		return nil
//	}
//	return v.Validate(obj)
//}
//
///**
// * 验证数据(开发者可以自由调用，但一般不推荐)
// * @param obj 需要验证的数据
// */
//func (v *Validator) Validate(obj interface{}) error {
//	if v.enable == false {
//		return nil
//	}
//
//	if err := v.v.Struct(obj); err != nil {
//		return errors.New(processErr(obj, err))
//	}
//
//	return nil
//}
//
///**
// * 判断是否过滤验证
// * @param ctx 请求上下文
// */
//func (v *Validator) isExclude(ctx *routing.Context) bool {
//	if v.enable == false {
//		return true
//	}
//
//	//匹配规则如下：
//	//1. 全等匹配：goods.List，表示只有goods.List方法的参数不会被验证
//	//2. 通配符匹配：goods.*，表示goods控制器下所有方法的参数不会被验证
//	//3. 正则匹配：goods.[A-Z0-9]+，表示goods控制器下符合该正则表达式的方法的参数不会被验证
//	objectName, methodName := invocation.ParseObjectNameAndMethodName(ctx)
//	for _, exclude := range v.excludes {
//		isObjectNameMatch := reflect.DeepEqual(objectName, exclude.objectName) || string_util.SymbolMatch(objectName, exclude.objectName) || string_util.RegexpMatch(objectName, exclude.objectName)
//		isMethodNameMatch := reflect.DeepEqual(methodName, exclude.mathodName) || string_util.SymbolMatch(methodName, exclude.mathodName) || string_util.RegexpMatch(methodName, exclude.mathodName)
//		if isObjectNameMatch && isMethodNameMatch {
//			return true
//		}
//	}
//	return false
//}
//
///**
// * 处理错误信息
// * @param obj 需要验证的数据
// * @param err 错误信息
// */
//func processErr(obj interface{}, err error) string {
//	if err == nil { //如果为nil 说明校验通过
//		return ""
//	}
//	invalid, ok := err.(*vld.InvalidValidationError) //如果是输入参数无效，则直接返回输入参数错误
//	if ok {
//		return "输入参数错误：" + invalid.Error()
//	}
//
//	errorList := make([]string, 0)
//	validationErrs := err.(vld.ValidationErrors) //断言是ValidationErrors
//	for _, validationErr := range validationErrs {
//		fieldName := validationErr.Field() //获取是哪个字段不符合格式
//		typeOf := reflect.TypeOf(obj)
//		// 如果是指针，获取其属性
//		if typeOf.Kind() == reflect.Ptr {
//			typeOf = typeOf.Elem()
//		}
//		field, _ := typeOf.FieldByName(fieldName) //通过反射获取filed
//		message := strings.TrimSpace(fmt.Sprintf("%s", field.Tag.Get("message")))
//		if len(message) > 0 {
//			errorList = append(errorList, message)
//		} else {
//			errorList = append(errorList, strings.TrimSpace(fmt.Sprintf("%s", validationErr)))
//		}
//	}
//	return strings.Join(errorList, ",")
//}
