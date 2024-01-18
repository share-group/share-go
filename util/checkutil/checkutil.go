package checkutil

import (
	"fmt"
	"regexp"
	"strings"
)

// 判断是否为整型数字
func IsInteger(str any) bool {
	_str := strings.TrimSpace(fmt.Sprintf("%v", str))
	return regexp.MustCompile(`^[0-9]+$`).MatchString(_str)
}
