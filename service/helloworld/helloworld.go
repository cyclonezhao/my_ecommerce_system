package helloworld

import (
	"context"
	"log"
	pb "my_ecommerce_system/proto/helloworld"
)

// 创建Server结构体，将SayHello方法注册为它的成员函数
type Server struct {
	pb.UnimplementedGreeterServer    // Server结构体继承了pb.UnimplementedGreeterServer结构体的所有方法
}

// 重写pb.GreeterServer.SayHello方法，实现业务逻辑
func (s *Server) SayHello(ctx context.Context, hellorequest *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received message from %v", hellorequest.GetName())
	return &pb.HelloReply{Message: "Hello " + hellorequest.GetName()}, nil
}
