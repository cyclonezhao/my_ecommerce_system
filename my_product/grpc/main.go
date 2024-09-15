package main

import (
	"context"
	"fmt"
	"log"
	pb "my_ecommerce_system/my_system_api/grpc/proto/helloworld"

	"flag"
	microservice "my_ecommerce_system/microservice"

	"time"

	my_client "my_ecommerce_system/pkg/client"

	"google.golang.org/grpc"
)

func main() {
	// 解析命令行参数
	flag.Parse()

	// 初始化etcd
	my_client.InitEtcdClient()

	// 通过“etcd 服务注册器”创建gRPC解析器
	resolver, err := microservice.NewEtcdResolver()
	if err != nil {
		log.Fatalf("Create etcd resolver error: %v", err)
	}

	// 创建与服务器的连接
	// target 原来是写死 "localhost:58090"，现在要改为etcd上的二级服务
	// 用resolver通过二级服务名称找到实例列表
	// 通过rr负载均衡访问实例列表
	pathServerName := microservice.GetPathServerName("my_system")
	conn, err := grpc.Dial(
		pathServerName,
		grpc.WithInsecure(),
		grpc.WithResolvers(resolver),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		log.Fatalf("无法连接到服务器：%v", err)
	}
	defer conn.Close()

	// 创建一个新的 gRPC 客户端
	client := pb.NewGreeterClient(conn)
	// 构建请求
	request := &pb.HelloRequest{
		Name: "John",
	}

	for i := 0; i < 10; i++ {
		// 创建2秒超时ctx
		ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
		// 调用 gRPC 方法
		response, err := client.SayHello(ctx, request)
		if err != nil {
			log.Fatalf("调用 gRPC 方法失败：%v", err)
		}
		// 打印响应
		fmt.Println(response.Message)
	}

	// 睡眠一会再结束
	log.Println("3秒后结束，客户端自动断开连接")
	time.Sleep(time.Second * 3)
}
