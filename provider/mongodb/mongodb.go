package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/share-group/share-go/provider/config"
	loggerFactory "github.com/share-group/share-go/provider/logger"
	"github.com/share-group/share-go/util/arrayutil"
	"github.com/share-group/share-go/util/jsonutil"
	"github.com/share-group/share-go/util/maputil"
	"github.com/share-group/share-go/util/stringutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
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
			logger.Fatal("only one default connection is allowed")
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
		logger.Fatal("%v", err)
	}
	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Fatal("%v", err)
	}
	logger.Info("mongodb connect %s success ...", uri)
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

func throwErrorIfNotNil(err error) {
	if err != nil {
		panic(err)
	}
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
// entity-数据实体; connectionName-连接名称
func EnsureIndex[T any](entity T, connectionName ...string) {
	ctx := context.Background()
	typ := reflect.TypeOf(entity)
	connection := GetInstance(connectionName...)
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
		throwErrorIfNotNil(err)

		key_s, _ := json.Marshal(indexModel.Keys)
		option_s, _ := json.Marshal(indexModel.Options)
		logger.Info("collection [%s] create index: %v, index options: %v", collection, string(key_s), jsonutil.RemoveNullValues(string(option_s)))
	}
}

// 反序列化mongodb的数据到指定是数据实体
//
// ctx-上下文; cursor-mongodb返回的游标; entity-数据实体
func DecodeList[T any](ctx context.Context, cursor *mongo.Cursor, entity T) []*T {
	defer cursor.Close(ctx)
	result := make([]*T, 0)
	for cursor.Next(ctx) {
		cursor.Decode(&entity)
		result = append(result, &entity)
	}
	return result
}

// 反序列化mongodb的数据到指定是数据实体
//
// ctx-上下文; cursor-mongodb返回的游标; entity-数据实体
func DecodeOne[T any](ctx context.Context, cursor *mongo.Cursor, entity T) *T {
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		cursor.Decode(&entity)
		return &entity
	}
	return nil
}

// 插入单条数据
//
// entity-数据实体
func (m *mongodb) InsertOne(entity any) primitive.ObjectID {
	return arrayutil.First(m.InsertMany(entity))
}

// 插入多条数据
//
// entity-数据实体(可以传数组或者传无限个单条)
func (m *mongodb) InsertMany(entity ...any) []primitive.ObjectID {
	ctx := context.Background()
	classType := reflect.TypeOf(arrayutil.First(entity))
	ignoreColumn := []string{"id", "createTime", "updateTime"}
	c := m.DB.Collection(strings.Split(fmt.Sprintf("%v", classType), ".")[1])
	documents := make([]any, 0)
	for _, item := range entity {
		document := make(bson.D, 0)
		pValue := reflect.ValueOf(item)
		for i := 0; i < classType.NumField(); i++ {
			fieldName := stringutil.FirstLowerCase(classType.Field(i).Name)
			if arrayutil.Contains(ignoreColumn, fieldName) {
				continue
			}
			fieldValue := pValue.Field(i).Interface()
			document = append(document, bson.E{Key: fieldName, Value: fieldValue})
		}

		now := time.Now().UnixMilli()
		document = append(document, bson.E{Key: "createTime", Value: now})
		document = append(document, bson.E{Key: "updateTime", Value: now})
		documents = append(documents, document)
	}

	result, err := c.InsertMany(ctx, documents)
	throwErrorIfNotNil(err)
	insertedIDs := make([]primitive.ObjectID, 0)
	for _, insertedID := range result.InsertedIDs {
		insertedIDs = append(insertedIDs, insertedID.(primitive.ObjectID))
	}
	return insertedIDs
}

// 更新单条数据
//
// entity-数据实体; id-主键id; update-需要更新的数据; opts-数据更新选项
func (m *mongodb) UpdateOne(entity any, id string, update bson.D, opts ...*options.UpdateOptions) *mongo.UpdateResult {
	objectID, err := primitive.ObjectIDFromHex(id)
	throwErrorIfNotNil(err)
	return m.UpdateMany(entity, bson.D{{"_id", objectID}}, update, opts...)
}

// 更新多条数据
//
// entity-数据实体; query-查询条件; update-需要更新的数据; opts-数据更新选项
func (m *mongodb) UpdateMany(entity any, query, update bson.D, opts ...*options.UpdateOptions) *mongo.UpdateResult {
	return m.doUpdateMany(entity, query, update, "$set", opts...)
}

// 单条数据的字段自增
//
// entity-数据实体; id-主键id; update-需要自增的数据; opts-数据更新选项
func (m *mongodb) IncOne(entity any, id string, update bson.D, opts ...*options.UpdateOptions) *mongo.UpdateResult {
	objectID, err := primitive.ObjectIDFromHex(id)
	throwErrorIfNotNil(err)
	return m.doUpdateMany(entity, bson.D{{"_id", objectID}}, update, "$inc", opts...)
}

// 更新多条数据
//
// entity-数据实体; query-查询条件; update-需要更新的数据; operator-操作符; opts-数据更新选项
func (m *mongodb) doUpdateMany(entity any, query, update bson.D, operator string, opts ...*options.UpdateOptions) *mongo.UpdateResult {
	ctx := context.Background()
	classType := reflect.TypeOf(entity)
	c := m.DB.Collection(strings.Split(fmt.Sprintf("%v", classType), ".")[1])
	updateResult, err := c.UpdateMany(ctx, query, bson.D{
		{operator, update},
		{"$set", bson.D{{Key: "updateTime", Value: time.Now().UnixMilli()}}},
	}, opts...)
	throwErrorIfNotNil(err)
	return updateResult
}

// 删除单条数据
//
// entity-数据实体; query-查询条件; opts-数据删除选项
func (m *mongodb) DeleteOne(entity any, query bson.D, opts ...*options.DeleteOptions) *mongo.DeleteResult {
	ctx := context.Background()
	classType := reflect.TypeOf(entity)
	c := m.DB.Collection(strings.Split(fmt.Sprintf("%v", classType), ".")[1])
	updateResult, err := c.DeleteOne(ctx, query, opts...)
	throwErrorIfNotNil(err)
	return updateResult
}

// 删除多条数据
//
// entity-数据实体; query-查询条件; opts-数据删除选项
func (m *mongodb) DeleteMany(entity any, query bson.D, opts ...*options.DeleteOptions) *mongo.DeleteResult {
	ctx := context.Background()
	classType := reflect.TypeOf(entity)
	c := m.DB.Collection(strings.Split(fmt.Sprintf("%v", classType), ".")[1])
	updateResult, err := c.DeleteMany(ctx, query, opts...)
	throwErrorIfNotNil(err)
	return updateResult
}

// 查询多条数据
//
// query-查询条件; entity-数据实体; opts-查询选项
func (m *mongodb) Find(query bson.D, entity any, opts ...*options.FindOptions) (context.Context, *mongo.Cursor) {
	ctx := context.Background()
	classType := reflect.TypeOf(entity)
	c := m.DB.Collection(strings.Split(fmt.Sprintf("%v", classType), ".")[1])
	cursor, err := c.Find(ctx, query, opts...)
	throwErrorIfNotNil(err)
	return ctx, cursor
}

// 查询单条数据
//
// query-查询条件; entity-数据实体
func (m *mongodb) FindOne(query bson.D, entity any) (context.Context, *mongo.Cursor) {
	opts := &options.FindOptions{}
	opts.SetLimit(1)
	return m.Find(query, entity, opts)
}

// 游标翻页
//
// query-查询条件; cursor-游标; pageSize-分页大小; sort-排序方式; entity-数据实体
func (m *mongodb) PaginationByCursor(query bson.D, cursor *string, pageSize int64, sort bson.D, entity any) (context.Context, *mongo.Cursor) {
	opts := &options.FindOptions{}
	opts.SetLimit(pageSize)
	if len(*cursor) > 0 {
		objectID, _ := primitive.ObjectIDFromHex(*cursor)
		query = bson.D{{"_id", bson.D{{"$lt", objectID}}}}
	}

	finalSort := bson.D{{"_id", -1}}
	if len(sort) > 0 {
		for _, s := range sort {
			finalSort = append(finalSort, s)
		}
	}
	opts.SetSort(finalSort)
	return m.Find(query, entity, opts)
}

// 页码翻页
//
// query-查询条件; page-当前页码; pageSize-分页大小; sort-排序方式; entity-数据实体
func (m *mongodb) PaginationByPage(query bson.D, page, pageSize int64, sort bson.D, entity any) (context.Context, *mongo.Cursor) {
	opts := &options.FindOptions{}
	opts.SetLimit(pageSize)
	opts.SetSkip((page - 1) * pageSize)
	if len(sort) <= 0 {
		sort = bson.D{{"_id", -1}}
	}
	opts.SetSort(sort)
	return m.Find(query, entity, opts)
}
