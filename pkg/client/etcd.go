package client

import (
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdClientWrapperStruct struct {
	EtcdClient *clientv3.Client
	EtcdUrl    string
	EtcdHost   string
	EtcdPort   int
}

func initEtcdClient() *EtcdClientWrapperStruct {
	protocal := "http"
	host := "127.0.0.1"
	port := 2379
	url := fmt.Sprintf("%s://%s:%d", protocal, host, port)

	/* 没有效果，以后有时间再看...
	// 设置自定义的 gRPC 重试策略
	backoffConfig := backoff.Config{
		BaseDelay:  1 * time.Second,  // 初始重试延迟
		Multiplier: 1.5,              // 重试延迟倍数
		MaxDelay:   10 * time.Second, // 最大重试延迟
	}

	// 使用自定义连接参数，配置重试策略
	connectParams := grpc.ConnectParams{
		Backoff:           backoffConfig,   // 应用自定义重试策略
		MinConnectTimeout: 5 * time.Second, // 最小连接超时
	}
	*/

	// 连接 etcd 客户端
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{url},
		DialTimeout: 5 * time.Second,
		/*DialOptions: []grpc.DialOption{
			grpc.WithConnectParams(connectParams), // 设置连接参数
			grpc.WithBlock(),                      // 阻塞等待连接
		},*/
	})
	if err != nil {
		log.Fatal(err)
	}

	// 连接成功，保存客户端实例
	return &EtcdClientWrapperStruct{etcdClient, url, host, port}
}
