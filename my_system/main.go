package main

import (
	"fmt"
	"log"
	"my_system/grpc/helloworld"
	"net"
	"os"
	"os/signal"
	"syscall"

	"flag"
	microservice "my_ecommerce_system/microservice"
	pb "my_ecommerce_system/my_system_api/grpc/proto/helloworld"
	my_client "my_ecommerce_system/pkg/client"

	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

var (
	grpcPort = flag.Int("port", 58090, "The server port")
)

type Config struct {
	DB      my_client.DbConfig    `yaml:"db"`
	Redis   my_client.RedisConfig `yaml:"redis"`
	Gateway struct {
		WriteList []string `yaml:"writeList"`
	} `yaml:"gateway"`
	Jwt struct {
		Expire int `yaml:"expire"`
	} `yaml:"jwt"`
}

var config Config

func updateConfigFn(rawConfig []byte) {
	// 将 YAML 字符串: rawConfig, 反序列化为结构体
	err := yaml.Unmarshal([]byte(rawConfig), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// 打印结构体内容
	fmt.Printf("%+v\n", config)
}

func main() {
	// 初始化etcd
	my_client.InitEtcdClient()
	// 拉取配置信息
	microservice.GetRawConfigFromConfigCenter("my_system", updateConfigFn)

	// 初始化数据库
	my_client.InitDB(&config.DB)

	// main退出后，关闭已经打开的第三方网络实体的客户端
	defer my_client.Close()

	// 解析命令行参数
	flag.Parse()

	// 绑定端口
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 初始化gRPC服务
	s := grpc.NewServer()
	// 将gRPC服务和自定义的业务逻辑注册到Greeter服务中
	pb.RegisterGreeterServer(s, &helloworld.Server{GRPCPort: grpcPort})
	log.Printf("serving gRPC on %v", lis.Addr())

	// 初始化ETCD注册器
	namingService, err := microservice.NewNamingService("my_system")
	if err != nil {
		log.Fatalf("failed to create NamingService: %v", err)
	}

	// 将本实例注册到ETCD
	err = namingService.AddEndpoint(microservice.Endpoint{
		Addr:    "localhost",
		Name:    "user",
		Port:    *grpcPort,
		Version: "1.0.0",
	})
	if err != nil {
		log.Fatalf("failed to reg etcd: %v", err)
	}

	// 启动RPC并监听。另起一个协程运行是为了避免  s.Serve(lis) 阻塞主协程
	// 后续主协程还要监听关闭信号
	log.Printf("server listening at %v", lis.Addr())
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
			namingService.DelAllEndpoint()
		}
	}()

	// 等待关闭信号
	quit := make(chan os.Signal)
	// syscall.SIGINT：通常是通过用户在终端按下 Ctrl+C 触发的中断信号
	// syscall.SIGTERM：是操作系统发送给程序的终止信号，可以通过 kill <pid> 命令发送 SIGTERM 信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Printf("Shutdown Server ... \r\n")
	// 停止grpc服务
	s.GracefulStop()
	// 删除etcd注册信息
	namingService.DelAllEndpoint()
	fmt.Printf("Graceful Shutdown Server success\r\n")
}
