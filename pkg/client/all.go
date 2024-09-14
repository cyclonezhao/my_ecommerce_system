package client

import (
	"database/sql"
	"fmt"
	"io"

	"github.com/go-redis/redis/v8"
)

var closerMap = make(map[string]io.Closer)

var RedisClient *redis.Client
var EtcdClientWrapper *EtcdClientWrapperStruct
var DB *sql.DB

// 关闭所有初始化的客户端
func Close() {
	fmt.Println("准备关闭客户端")
	for cli, closer := range closerMap {
		fmt.Printf("关闭%s\n", cli)
		closer.Close()
	}
	fmt.Println("客户端已全部关闭")
}

// 初始化etcd
func InitEtcdClient() {
	EtcdClientWrapper = initEtcdClient()
	closerMap["etcd"] = EtcdClientWrapper.EtcdClient
}

// 初始化DB，比如MySQL
func InitDB(config *DbConfig) {
	DB = initDB(config)
	closerMap["db"] = DB
}

// 初始化Redis
func InitRedis(config *RedisConfig) {
	RedisClient = initRedis(config)
	closerMap["redis"] = RedisClient
}
