package client

import (
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DB       int    `yaml:"db"`
	Password string `yaml:"password"`
}

func initRedis(config *RedisConfig) *redis.Client {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		DB:       config.DB,
		Password: config.Password,
	})
	log.Println("Redis Client 初始化成功")
	return RedisClient
}
