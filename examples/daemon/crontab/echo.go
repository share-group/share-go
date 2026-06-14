package crontab

import (
	"log"
	"time"

	json "github.com/bytedance/sonic"
	entity "github.com/share-group/share-go/examples/daemon/entity/account"
	"github.com/share-group/share-go/provider/db/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var dashboard = mongodb.GetInstance(entity.Role{}, "dashboard")

func init() {
	go func() {
		for {
			time.Sleep(time.Second)
			b, _ := json.Marshal(dashboard.Find(bson.D{}))
			log.Println("hello world", string(b))
		}
	}()
}
