package util

import (
	"fmt"
	"net/url"
	"strings"
)

type httpUtil struct{}

var HttpUtil = newHttpUtil()

func newHttpUtil() *httpUtil {
	return &httpUtil{}
}

// 解析 querystring
//
// urlString
func (s *httpUtil) ParseQueryString(urlString string) string {
	queryStringMap := make(map[string]any)
	parsedURL, _ := url.Parse(urlString)
	for k, v := range parsedURL.Query() {
		var newValue any
		if len(v) == 1 {
			newValue = strings.TrimSpace(fmt.Sprintf("%v", v[0]))
		} else {
			newValue = v
		}
		queryStringMap[k] = newValue
	}
	return StringUtil.JSON(queryStringMap)
}
