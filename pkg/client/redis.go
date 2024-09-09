package client

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"my_ecommerce_system/pkg/config"
)

var RedisClient *redis.Client

func InitRedis(){
	RedisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", config.AppConfig.Redis.Host, config.AppConfig.Redis.Port),
		DB:   config.AppConfig.Redis.DB,
		Password: config.AppConfig.Redis.Password,
	})
	log.Println("Redis Client 初始化成功")
}
