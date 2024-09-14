package main

import (
	"fmt"
	"log"
	"my_system/grpc/helloworld"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"flag"
	microservice "my_ecommerce_system/microservice"
	pb "my_ecommerce_system/my_system_api/grpc/proto/helloworld"
	my_client "my_ecommerce_system/pkg/client"

	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

var (
	grpcPort = flag.Int("port", 8081, "The server port")
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

	// 创建 cmux 实例，用于在同一个端口同时监听处理gRPC和http请求
	m := cmux.New(lis)
	var cmuxClosed = false
	// 匹配 gRPC 请求
	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	// 匹配 HTTP 请求
	httpL := m.Match(cmux.HTTP1Fast())

	// 初始化gRPC服务
	grpcServer := grpc.NewServer()
	// 将gRPC服务和自定义的业务逻辑注册到Greeter服务中
	pb.RegisterGreeterServer(grpcServer, &helloworld.Server{GRPCPort: grpcPort})
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

	// 启动RPC并监听。另起一个协程运行是为了避免  s.Serve(grpcL) 阻塞主协程
	// 后续主协程还要监听关闭信号
	log.Printf("server listening at %v", grpcL.Addr())
	go func() {
		if err := grpcServer.Serve(grpcL); err != nil {
			if cmuxClosed {
				return
			}
			log.Fatalf("failed to serve gRPC: %v", err)
			namingService.DelAllEndpoint()
		}
	}()

	// 启动HTTP服务
	engine := gin.Default()
	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	httpServer := &http.Server{
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		log.Printf("HTTP server listening at %v", grpcL.Addr())
		if err := httpServer.Serve(httpL); err != nil {
			if cmuxClosed {
				return
			}
			log.Fatalf("Failed to serve HTTP: %v", err)
		}
	}()

	go func() {
		if err := m.Serve(); err != nil {
			if cmuxClosed {
				return
			}
			log.Fatalf("Failed to serve c: %v", err)
		}
	}()

	// 等待关闭信号
	quit := make(chan os.Signal, 1)
	// syscall.SIGINT：通常是通过用户在终端按下 Ctrl+C 触发的中断信号
	// syscall.SIGTERM：是操作系统发送给程序的终止信号，可以通过 kill <pid> 命令发送 SIGTERM 信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutdown Server ...")
	// cmux必须关掉，不然下面的grpcServer.GracefulStop()会阻塞
	// 原因猜测是cmux维持着对底层连接的控制，导致 grpcServer.GracefulStop() 持续等待，阻塞停止流程
	cmuxClosed = true
	m.Close()
	// 停止grpc服务
	grpcServer.GracefulStop()
	fmt.Println("gRPC服务已停止")
	// 删除etcd注册信息
	namingService.DelAllEndpoint()
	fmt.Println("Graceful Shutdown Server success.")
}
