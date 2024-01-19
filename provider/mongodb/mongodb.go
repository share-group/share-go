package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/share-group/share-go/provider/config"
	loggerFactory "github.com/share-group/share-go/provider/logger"
	"github.com/share-group/share-go/util/jsonutil"
	"github.com/share-group/share-go/util/maputil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var connectionMap sync.Map
var logger = loggerFactory.GetLogger()

type mongodb struct {
	DB *mongo.Database
}

func init() {
	for _, conf := range config.GetList("data.mongodb") {
		name := fmt.Sprintf("%v", maputil.GetValueFromMap(conf.(map[string]any), "name", "default"))
		uri := fmt.Sprintf("%v", maputil.GetValueFromMap(conf.(map[string]any), "uri", ""))

		if _, ok := connectionMap.Load(name); ok {
			logger.DPanic("only one default connection is allowed")
			os.Exit(1)
		}
		connectionMap.Store(name, newMongodb(name, uri))
	}
}

func newMongodb(name string, uri string) *mongodb {
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

func GetInstance(connectionName ...string) *mongodb {
	if len(connectionName) <= 0 || len(connectionName[0]) <= 0 {
		connectionName = append(connectionName, "default")
	}

	connection, ok := connectionMap.Load(connectionName[0])
	if !ok {
		panic(fmt.Sprintf("unable to get [%s] connection", connectionName[0]))
	}
	return connection.(*mongodb)
}

// 创建索引
//
// entity-数据实体;connectionName-连接名称
func EnsureIndex[T any](entity T, connectionName ...string) {
	ctx := context.Background()
	typ := reflect.TypeOf(entity)
	connection := GetInstance(connectionName...)
	connection.DB.Drop(ctx)
	collection := typ.Name()
	for i := 0; i < typ.NumField(); i++ {
		indexJSON := strings.TrimSpace(typ.Field(i).Tag.Get("index"))
		if len(indexJSON) <= 0 {
			continue
		}

		indexMap := make(map[string]any)
		indexJSON = strings.ReplaceAll(indexJSON, "'", "\"")
		json.Unmarshal([]byte(indexJSON), &indexMap)
		if len(indexMap) <= 0 {
			// 忽略解析失败
			continue
		}

		keys := make(bson.D, 0)
		regexpPattern, _ := regexp.Compile("\\s+")
		indexJSON = regexpPattern.ReplaceAllString(indexJSON, "")
		keyString := indexJSON[strings.Index(indexJSON, `"keys":{`)+8 : strings.Index(indexJSON, `}`)]
		for _, key := range strings.Split(keyString, ",") {
			arr := strings.Split(key, ":")
			k := strings.ReplaceAll(strings.TrimSpace(arr[0]), `"`, "")
			v, _ := strconv.Atoi(strings.TrimSpace(arr[1]))
			keys = append(keys, bson.E{Key: k, Value: v})
		}

		indexName := fmt.Sprintf("%v", maputil.GetValueFromMap(indexMap, "name", ""))
		sparse, _ := strconv.ParseBool(fmt.Sprintf("%v", maputil.GetValueFromMap(indexMap, "sparse", false)))
		unique, _ := strconv.ParseBool(fmt.Sprintf("%v", maputil.GetValueFromMap(indexMap, "unique", false)))
		indexModel := mongo.IndexModel{Keys: keys}
		indexModel.Options = options.Index().SetName(indexName).SetBackground(true)
		if unique {
			indexModel.Options.SetUnique(unique)
		}
		if sparse {
			indexModel.Options.SetSparse(sparse)
		}
		_, err := connection.DB.Collection(collection).Indexes().CreateOne(ctx, indexModel)
		if err != nil {
			panic(err)
		}

		key_s, _ := json.Marshal(indexModel.Keys)
		option_s, _ := json.Marshal(indexModel.Options)
		logger.Info(fmt.Sprintf("collection [%s] create index: %v, index options: %v", collection, string(key_s), jsonutil.RemoveNullValues(string(option_s))))
	}
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
	err := cursor.All(ctx, &slice)
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
