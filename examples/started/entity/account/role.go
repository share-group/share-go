package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Role struct {
	Id         primitive.ObjectID `bson:"_id"`
	Name       string             `bson:"name"`
	Authority  any                `bson:"authority"`
	CreateTime int64              `bson:"createTime"`
}
