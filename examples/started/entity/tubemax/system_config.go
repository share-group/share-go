package entity

import (
	"github.com/share-group/share-go/provider/db"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type SystemConfig struct {
	Id     bson.ObjectID `bson:"_id"`
	Module string        `bson:"module" index:"{'name':'moduleAndKey','keys':{'module':-1,'key':-1},'unique':true}"`
	Key    string        `bson:"key"`
	Value  any           `bson:"value"`
	Desc   string        `bson:"desc"`
}

func init() {
	db.RegisterRepository(SystemConfig{}, "tubemax")
}
