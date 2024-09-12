package main

import (
	"context"
	"fmt"
	"log"
	pb "my_ecommerce_system/my_system_api/grpc/proto/helloworld"

	"google.golang.org/grpc"
)

func main() {
	// 创建与服务器的连接
	conn, err := grpc.Dial("localhost:58090", grpc.WithInsecure())
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
	// 调用 gRPC 方法
	response, err := client.SayHello(context.Background(), request)
	if err != nil {
		log.Fatalf("调用 gRPC 方法失败：%v", err)
	}
	// 打印响应
	fmt.Println(response.Message)
}
