package main

import (
	"fmt"
	"log"
	pb "my_ecommerce_system/my_system_api/grpc/proto/helloworld"
	"my_system/grpc/helloworld"
	"net"

	"google.golang.org/grpc"
)

func main() {
	// 绑定地址和端口
	grpcAddress := "0.0.0.0"
	grpcPort := 58090
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", grpcAddress, grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 初始化gRPC服务
	s := grpc.NewServer()
	// 将gRPC服务和自定义的业务逻辑注册到Greeter服务中
	pb.RegisterGreeterServer(s, &helloworld.Server{})
	log.Printf("serving gRPC on %v", lis.Addr())
	// 将gRPC服务绑定在上面创建的tcp端口上，并开启监听
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("启动gRPC服务失败(%v)", err)
	}
}
