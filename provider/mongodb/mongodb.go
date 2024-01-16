package mongodb

import (
	"context"
	"fmt"
	"github.com/share-group/share-go/provider/config"
	loggerFactory "github.com/share-group/share-go/provider/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"reflect"
	"strings"
)

var logger = loggerFactory.GetLogger()

type mongodb struct {
	DB *mongo.Database
}

var Mongodb = NewMongodb(config.GetString("data.mongodb.uri"))

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
		logger.Fatal(fmt.Sprintf("%v", err))
		os.Exit(1)
	}
	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Fatal(fmt.Sprintf("%v", err))
		os.Exit(1)
	}
	logger.Info(fmt.Sprintf("mongodb connect %s success ...", uri))
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

// 查询数据
//
// query-查询条件; entity-数据实体
func (m *mongodb) Find(query bson.D, entity any, opts ...*options.FindOptions) any {
	ctx := context.Background()
	classType := reflect.TypeOf(entity)
	c := m.DB.Collection(strings.Split(fmt.Sprintf("%v", classType), ".")[1])
	cursor, _ := c.Find(ctx, query, opts...)
	defer cursor.Close(ctx)
	slice := reflect.MakeSlice(reflect.SliceOf(classType), 1, 1).Interface()
	err := cursor.All(context.Background(), &slice)
	if err != nil {
		logger.DPanic(err.Error())
	}
	return slice
}

// 游标翻页
//
// query-查询条件; cursor-游标; pageSize-分页大小; sort-排序方式; entity-数据实体
func (m *mongodb) PaginationByCursor(query bson.D, cursor *string, pageSize int64, sort bson.D, entity any) any {
	opts := &options.FindOptions{}
	opts.SetLimit(pageSize)
	if len(*cursor) > 0 {
		objectID, _ := primitive.ObjectIDFromHex(*cursor)
		query = bson.D{{"_id", bson.D{{"$lt", objectID}}}}
	}
	if len(sort) <= 0 {
		sort = bson.D{{"_id", -1}}
	}
	opts.SetSort(sort)
	return m.Find(query, entity, opts)
}

// 页码翻页
//
// query-查询条件; page-当前页码; pageSize-分页大小; sort-排序方式; entity-数据实体
func (m *mongodb) PaginationByPage(query bson.D, page, pageSize int64, sort bson.D, entity any) any {
	opts := &options.FindOptions{}
	opts.SetLimit(pageSize)
	opts.SetSkip((page - 1) * pageSize)
	if len(sort) <= 0 {
		sort = bson.D{{"_id", -1}}
	}
	opts.SetSort(sort)
	return m.Find(query, entity, opts)
}
