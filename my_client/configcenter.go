package client

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	CONFIG_PREFIX = "my_ecommerce_system/config/"
)

// 从配置中心拉取原始配置信息
func GetRawConfigFromConfigCenter(appName string, updateConfigFn func(yamlStr []byte)) {
	// 连接 etcd 客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// 从 etcd 获取配置信息
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	key := CONFIG_PREFIX + appName
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.Kvs) > 0 {
		fmt.Printf("配置加载成功: %s\n", key)
		updateConfigFn(resp.Kvs[0].Value)
	}

	// 开启监听配置变化
	go watchConfig(cli, key, updateConfigFn)
}

func watchConfig(cli *clientv3.Client, key string, updateConfigFn func(yamlStr []byte)) {
	rch := cli.Watch(context.Background(), key)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("配置发生变化: %s\n", key)
			updateConfigFn(ev.Kv.Value)
		}
	}
}
