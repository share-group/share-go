package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/share-group/share-go/provider/config"
	loggerFactory "github.com/share-group/share-go/provider/logger"
	"github.com/share-group/share-go/provider/mongodb"
	"github.com/share-group/share-go/util/systemutil"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"strings"
	"time"
)

var logger = loggerFactory.GetLogger()
var loggingEnable = config.GetBool("server.logging.enable")
var loggingPretty = config.GetBool("server.logging.pretty")
var loggingMongodb = mongodb.GetInstance("logging")

func PrintRequestLog(c echo.Context) {
	// 正式环境不打印请求日志
	if !loggingEnable || reflect.DeepEqual(systemutil.Env(), "prod") || reflect.DeepEqual(systemutil.Env(), "production") {
		return
	}

	req := c.Get("request")
	requestBytes := make([]byte, 0)
	if req != nil {
		requestBytes = req.([]byte)
	}
	if loggingPretty {
		request := make(map[string]any)
		json.Unmarshal(requestBytes, &request)
		requestBytes, _ = json.MarshalIndent(request, "", "    ")
		c.Set("request", requestBytes)
	}

	logger.Info("request %v, method: %v, data: %v, size: %v Byte", c.Request().URL.Path, c.Request().Method, string(requestBytes), len(requestBytes))
}

func SaveStringRequestLog(c echo.Context) {
	req := c.Get("request")
	requestTime := c.Get("requestTime")
	requestBytes := make([]byte, 0)
	if req != nil {
		requestBytes = req.([]byte)
	}
	responseBytes := c.Get("response").([]byte)
	exec := time.Since(requestTime.(time.Time))

	if loggingMongodb != nil {
		logEntity := bson.D{
			bson.E{Key: "machine", Value: systemutil.GetHostName()},
			bson.E{Key: "url", Value: c.Request().URL.Path},
			bson.E{Key: "originUrl", Value: c.Request().RequestURI},
			bson.E{Key: "method", Value: c.Request().Method},
			bson.E{Key: "ip", Value: c.RealIP()},
		}

		headers := bson.D{}
		for header, name := range c.Request().Header {
			key := strings.TrimSpace(header)
			value := strings.TrimSpace(strings.Join(name, ";"))
			if len(key) <= 0 || len(value) <= 0 {
				continue
			}
			headers = append(headers, bson.E{Key: key, Value: value})
		}

		logEntity = append(logEntity, bson.E{Key: "headers", Value: headers})
		logEntity = append(logEntity, bson.E{Key: "request", Value: string(requestBytes)})
		logEntity = append(logEntity, bson.E{Key: "response", Value: string(responseBytes)})
		logEntity = append(logEntity, bson.E{Key: "status", Value: c.Response().Status})
		logEntity = append(logEntity, bson.E{Key: "duration", Value: exec.String()})
		logEntity = append(logEntity, bson.E{Key: "requestTime", Value: c.Get("requestTime").(time.Time).UnixMilli()})
		logEntity = append(logEntity, bson.E{Key: "responseTime", Value: time.Now().UnixMilli()})
		go loggingMongodb.DB.Collection(fmt.Sprintf("Log_%s", time.Now().Format("200601"))).InsertOne(context.Background(), logEntity)
	}

	// 测试环境打印详细点，正式环境打印简单点
	if loggingEnable {
		if reflect.DeepEqual(systemutil.Env(), "production") {
			logger.Info("%v %v %v %v", c.Response().Status, c.Request().Method, c.Request().URL.Path, exec)
		} else {
			if loggingPretty {
				c.Set("response", responseBytes)
			}
			logger.Info("response %v %v, data: %v, size: %v Byte, exec: %v", c.Request().URL.Path, c.Response().Status, string(responseBytes), len(responseBytes), exec)
		}
	}
}

func SaveJSONRequestLog(c echo.Context) {
	req := c.Get("request")
	requestTime := c.Get("requestTime")
	request := make(map[string]any)
	response := make(map[string]any)
	requestBytes := make([]byte, 0)
	if req != nil {
		requestBytes = req.([]byte)
	}
	responseBytes := c.Get("response").([]byte)
	json.Unmarshal(requestBytes, &request)
	json.Unmarshal(responseBytes, &response)
	exec := time.Since(requestTime.(time.Time))

	if loggingMongodb != nil {
		logEntity := bson.D{
			bson.E{Key: "machine", Value: systemutil.GetHostName()},
			bson.E{Key: "url", Value: c.Request().URL.Path},
			bson.E{Key: "originUrl", Value: c.Request().RequestURI},
			bson.E{Key: "method", Value: c.Request().Method},
			bson.E{Key: "ip", Value: c.RealIP()},
		}

		headers := bson.D{}
		for header, name := range c.Request().Header {
			key := strings.TrimSpace(header)
			value := strings.TrimSpace(strings.Join(name, ";"))
			if len(key) <= 0 || len(value) <= 0 {
				continue
			}
			headers = append(headers, bson.E{Key: key, Value: value})
		}

		logEntity = append(logEntity, bson.E{Key: "headers", Value: headers})
		logEntity = append(logEntity, bson.E{Key: "request", Value: request})
		logEntity = append(logEntity, bson.E{Key: "response", Value: response})
		logEntity = append(logEntity, bson.E{Key: "status", Value: c.Response().Status})
		logEntity = append(logEntity, bson.E{Key: "duration", Value: exec.String()})
		logEntity = append(logEntity, bson.E{Key: "requestTime", Value: c.Get("requestTime").(time.Time).UnixMilli()})
		logEntity = append(logEntity, bson.E{Key: "responseTime", Value: time.Now().UnixMilli()})
		go loggingMongodb.DB.Collection(fmt.Sprintf("Log_%s", time.Now().Format("200601"))).InsertOne(context.Background(), logEntity)
	}

	// 测试环境打印详细点，正式环境打印简单点
	if loggingEnable {
		if reflect.DeepEqual(systemutil.Env(), "production") {
			logger.Info("%v %v %v %v", c.Response().Status, c.Request().Method, c.Request().URL.Path, exec)
		} else {
			if loggingPretty {
				json.Unmarshal(responseBytes, &response)
				responseBytes, _ = json.MarshalIndent(response, "", "    ")
				c.Set("response", responseBytes)
			}
			logger.Info("response %v %v, data: %v, size: %v Byte, exec: %v", c.Request().URL.Path, c.Response().Status, string(responseBytes), len(responseBytes), exec)
		}
	}
}
