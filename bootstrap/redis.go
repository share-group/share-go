package bootstrap

import (
	"github.com/go-redis/redis"
	"log"
)

type redisObj struct {
	client *redis.Client
}

var Redis = NewRedis(Config.GetStringValue("redis.host"), Config.GetStringValue("redis.password"), Config.GetIntegerValue("redis.db"))

var redisClient *redis.Client

func NewRedis(host, password string, db int) *redisObj {
	redisClient = redis.NewClient(&redis.Options{Addr: host, Password: password, DB: db})
	ping, err := redisClient.Ping().Result()
	log.Println("ping: ", ping)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Redis connected:", host)
	return &redisObj{client: redisClient}
}

func GetRedisClient() *redis.Client {
	return redisClient
}
