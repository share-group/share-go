package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/share-group/share-go/provider/config"
	loggerFactory "github.com/share-group/share-go/provider/logger"
	"github.com/share-group/share-go/util/maputil"
	"github.com/share-group/share-go/util/systemutil"
	"os"
	"strconv"
	"sync"
)

var connectionMap sync.Map
var logger = loggerFactory.GetLogger()

type redisObj struct {
	Client *redis.Client
}

var redisClient *redis.Client

func init() {
	for _, conf := range config.GetList("data.redis") {
		name := fmt.Sprintf("%v", maputil.GetValueFromMap(conf.(map[string]any), "name", "default"))
		host := fmt.Sprintf("%v", maputil.GetValueFromMap(conf.(map[string]any), "host", ""))
		password := fmt.Sprintf("%v", maputil.GetValueFromMap(conf.(map[string]any), "password", ""))
		db, _ := strconv.Atoi(fmt.Sprintf("%v", maputil.GetValueFromMap(conf.(map[string]any), "host", "0")))

		if _, ok := connectionMap.Load(name); ok {
			logger.DPanic("only one default connection is allowed")
			os.Exit(1)
		}
		connectionMap.Store(name, newRedis(host, password, db))
	}
}

func newRedis(host, password string, db int) *redisObj {
	redisClient = redis.NewClient(&redis.Options{Addr: host, Password: password, DB: db})
	ping, err := redisClient.Ping().Result()
	logger.Info(fmt.Sprintf("redis connect %v %s ...", host, systemutil.If(ping == "PONG", "success", "fail")))
	if err != nil {
		logger.Fatal(fmt.Sprintf("%v", err))
		os.Exit(1)
	}
	return &redisObj{Client: redisClient}
}

func GetInstance(connectionName ...string) *redisObj {
	if len(connectionName) <= 0 || len(connectionName[0]) <= 0 {
		connectionName = append(connectionName, "default")
	}

	connection, ok := connectionMap.Load(connectionName[0])
	if !ok {
		panic(fmt.Sprintf("unable to get [%s] connection", connectionName[0]))
	}
	return connection.(*redisObj)
}
