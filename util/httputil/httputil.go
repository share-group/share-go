package httputil

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	json "github.com/bytedance/sonic"

	"github.com/pkg/errors"
	loggerFactory "github.com/share-group/share-go/provider/logger"
)

var logger = loggerFactory.GetLogger()

// 发送 post 请求
func Post(url string, headers map[string]string, body []byte, target interface{}) error {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		logger.Error("http.NewRequest error: %v", err)
		return err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("send http post request error: %v", err)
		return err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("read data error: %v", err)
		return err
	}

	err = json.Unmarshal(b, &target)
	if err != nil {
		logger.Error("send http post: %v %v", url, string(body))
		logger.Error("send http post response: %v", string(b))
		logger.Error("send http post response decode error: %v", err)
		return errors.Wrap(err, fmt.Sprintf("response %s", string(b)))
	}
	return nil
}

// 发送 get 请求
func Get(url string, headers map[string]string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		message := fmt.Sprintf("http.NewRequest error: %v", err)
		logger.Error(message)
		b, _ := json.Marshal(map[string]any{"code": 1, "errorCode": 1, "message": message, "errorMessage": message})
		return b
	}

	// 设置请求头
	if len(headers) > 0 {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(req)
	if err != nil {
		message := fmt.Sprintf("send http get request error: %v", err)
		logger.Error(message)
		b, _ := json.Marshal(map[string]any{"code": 1, "errorCode": 1, "message": message, "errorMessage": message})
		return b
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		message := fmt.Sprintf("read data error: %v", err)
		logger.Error(message)
		b, _ = json.Marshal(map[string]any{"code": 1, "errorCode": 1, "message": message, "errorMessage": message})
		return b
	}

	if resp.StatusCode != http.StatusOK {
		message := fmt.Sprintf("send http get %s request error, statusCode: %d", url, resp.StatusCode)
		logger.Error(message)
		if strings.Contains(url, "ip.me") {
			body, _ := json.Marshal(map[string]any{"msg_type": "text", "content": map[string]string{"text": message}})
			Post("https://open.feishu.cn/open-apis/bot/v2/hook/4882de36-de41-4fcc-a4c4-a1b186dae668", make(map[string]string), body, make(map[string]any))
		}
		b, _ = json.Marshal(map[string]any{"code": resp.StatusCode, "errorCode": resp.StatusCode, "message": message, "errorMessage": message})
		return b
	}

	return b
}

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

	b, _ := json.Marshal(queryStringMap)
	return string(b)
}
