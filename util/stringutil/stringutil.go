package stringutil

import (
	"fmt"
	"strings"
)

// 首字母大写
func FirstUpperCase(str string) string {
	return fmt.Sprintf("%s%s", strings.ToUpper(str[:1]), str[1:])
}

// 首字母小写
func FirstLowerCase(str string) string {
	return fmt.Sprintf("%s%s", strings.ToLower(str[:1]), str[1:])
}
