package util

import (
	"fmt"
	"regexp"
	"strings"
)

type checkUtil struct{}

var CheckUtil = newCheckUtil()

func newCheckUtil() *checkUtil {
	return &checkUtil{}
}

// 判断是否为整型数字
func (c *checkUtil) IsInteger(str any) bool {
	_str := strings.TrimSpace(fmt.Sprintf("%v", str))
	return regexp.MustCompile(`^[0-9]+$`).MatchString(_str)
}
