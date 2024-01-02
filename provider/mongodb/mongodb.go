package mongodb

import (
	"context"
	"fmt"
	"github.com/share-group/share-go/bootstrap"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strings"
)

var logger = bootstrap.Logger.GetLogger()

// 查询数据
//
// query-查询条件; entity-数据实体
func Find(query bson.D, entity any, opts ...*options.FindOptions) any {
	ctx := context.Background()
	classType := reflect.TypeOf(entity)
	c := bootstrap.Mongodb.DB.Collection(strings.Split(fmt.Sprintf("%v", classType), ".")[1])
	cursor, _ := c.Find(ctx, query, opts...)
	defer cursor.Close(ctx)
	slice := reflect.MakeSlice(reflect.SliceOf(classType), 1, 1).Interface()
	err := cursor.All(context.Background(), &slice)
	if err != nil {
		logger.DPanic(err.Error())
	}
	return slice
}

func PaginationByCursor(query bson.D, cursor string, pageSize int64, sort bson.D, entity any) any {
	opts := &options.FindOptions{}
	opts.SetLimit(pageSize)
	if len(cursor) > 0 {
		objectID, _ := primitive.ObjectIDFromHex(cursor)
		query = bson.D{{"_id", bson.D{{"$lt", objectID}}}}
	}
	if len(sort) <= 0 {
		sort = bson.D{{"_id", -1}}
	}
	opts.SetSort(sort)
	return Find(query, entity, opts)
}

func PaginationByPage(query bson.D, page, pageSize int64, sort bson.D, entity any) any {
	opts := &options.FindOptions{}
	opts.SetLimit(pageSize)
	opts.SetSkip((page - 1) * pageSize)
	if len(sort) <= 0 {
		sort = bson.D{{"_id", -1}}
	}
	opts.SetSort(sort)
	return Find(query, entity, opts)
}
