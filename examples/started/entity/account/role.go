package entity

import (
	"github.com/share-group/share-go/provider/db"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Role struct {
	Id         bson.ObjectID `bson:"_id"`
	Name       string        `bson:"name" index:"{'name':'name','keys':{'name':-1},'unique':true}"`
	Authority  any           `bson:"authority"`
	CreateTime int64         `bson:"createTime"`
}

func init() {
	db.RegisterRepository(Role{}, "dashboard")
}
