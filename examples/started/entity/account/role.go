package entity

import (
	"github.com/share-group/share-go/provider/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	Id         primitive.ObjectID `bson:"_id"`
	Name       string             `bson:"name" index:"{'name':'name','keys':{'name':-1},'unique':true}"`
	Authority  any                `bson:"authority"`
	CreateTime int64              `bson:"createTime"`
}

func init() {
	mongodb.EnsureIndex(Role{}, "dashboard")
}
