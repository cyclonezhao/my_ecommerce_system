package helloworld

import (
	"context"
	"fmt"
	"log"
	pb "my_ecommerce_system/my_system_api/grpc/proto/helloworld"
)

// 创建Server结构体，将SayHello方法注册为它的成员函数
type Server struct {
	GRPCPort                      *int
	pb.UnimplementedGreeterServer // Server结构体继承了pb.UnimplementedGreeterServer结构体的所有方法
}

// 重写pb.GreeterServer.SayHello方法，实现业务逻辑
func (s *Server) SayHello(ctx context.Context, hellorequest *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received message from %v", hellorequest.GetName())
	return &pb.HelloReply{Message: fmt.Sprintf("%d says: Hello %s", *s.GRPCPort, hellorequest.GetName())}, nil
}
