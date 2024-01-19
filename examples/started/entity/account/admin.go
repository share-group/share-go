package entity

import (
	"github.com/share-group/share-go/provider/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Admin struct {
	Id            primitive.ObjectID `bson:"_id"`
	Email         string             `bson:"email" index:"{'name':'email','keys':{'email':-1,'roleId':-1,'createTime':-1,'lastLoginIp':-1},'unique':true}"`
	Username      string             `bson:"username"`
	Password      string             `bson:"password"`
	Status        string             `bson:"status"`
	RoleId        string             `bson:"roleId"`
	CreateTime    int64              `bson:"createTime"`
	LastLoginTime int64              `bson:"lastLoginTime"`
	LastLoginIp   string             `bson:"lastLoginIp"`
}

func init() {
	mongodb.EnsureIndex(Admin{}, "dashboard")
}
