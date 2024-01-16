package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/share-group/share-go/provider/config"
	loggerFactory "github.com/share-group/share-go/provider/logger"
	"github.com/share-group/share-go/util"
	"os"
)

var logger = loggerFactory.GetLogger()

type redisObj struct {
	Client *redis.Client
}

var Redis = NewRedis(config.GetString("data.redis.host"), config.GetString("data.redis.password"), config.GetInt("data.redis.db"))

var redisClient *redis.Client

func NewRedis(host, password string, db int) *redisObj {
	redisClient = redis.NewClient(&redis.Options{Addr: host, Password: password, DB: db})
	ping, err := redisClient.Ping().Result()
	logger.Info(fmt.Sprintf("redis connect %v %s ...", host, util.SystemUtil.If(ping == "PONG", "success", "fail")))
	if err != nil {
		logger.Fatal(fmt.Sprintf("%v", err))
		os.Exit(1)
	}
	return &redisObj{Client: redisClient}
}
