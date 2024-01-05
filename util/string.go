package util

import (
	"encoding/json"
	"fmt"
	"strings"
)

type stringUtil struct{}

var StringUtil = newStringUtil()

func newStringUtil() *stringUtil {
	return &stringUtil{}
}

// 转json
func (s *stringUtil) JSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// 首字母大写
func (s *stringUtil) FirstUpperCase(str string) string {
	return fmt.Sprintf("%s%s", strings.ToUpper(str[:1]), str[1:])
}

// 首字母小写
func (s *stringUtil) FirstLowerCase(str string) string {
	return fmt.Sprintf("%s%s", strings.ToLower(str[:1]), str[1:])
}
