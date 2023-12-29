package bootstrap

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
)

type mongodb struct {
	DB *mongo.Database
}

var Mongodb = NewMongodb(Config.GetStringValue("data.mongodb.uri"))

func NewMongodb(uri string) *mongodb {
	if len(uri) <= 0 {
		return nil
	}

	// 设置客户端连接配置
	ctx := context.Background()
	co := options.Client().ApplyURI(uri)
	// 连接到MongoDB
	client, err := mongo.Connect(ctx, co)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Println(fmt.Sprintf("mongodb connect %s success ...", uri))
	return &mongodb{DB: client.Database(dbName(uri))}
}

func dbName(uri string) string {
	lastIndexSlash := strings.LastIndex(uri, "/") + 1
	lastIndexFactor := strings.LastIndex(uri, "?")
	if lastIndexFactor > 0 {
		return strings.TrimSpace(uri[lastIndexSlash:lastIndexFactor])
	}
	return strings.TrimSpace(uri[lastIndexSlash:])
}
