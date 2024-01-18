package httputil

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// 解析 querystring
//
// urlString-地址
func ParseQueryString(urlString string) string {
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

	bytes, _ := json.Marshal(queryStringMap)
	return string(bytes)
}
