package main

import (
	"net/http"

	"flag"

	"time"

	my_client "my_ecommerce_system/pkg/client"
	"my_product/internal/config"
	"my_product/internal/httprouter"

	microservice "my_ecommerce_system/pkg/microservice"

	"github.com/gin-gonic/gin"
)

func main() {
	// 解析命令行参数
	flag.Parse()

	// 初始化etcd 等客户端
	my_client.InitEtcdClient()
	// 拉取配置信息
	microservice.GetRawConfigFromConfigCenter("my_product", &config.AppConfig)
	my_client.InitXORM(&config.AppConfig.DB)
	defer my_client.Close()

	// 启动HTTP服务
	engine := gin.Default()
	httprouter.Route(engine)

	handler := engine

	httpServer := &http.Server{
		Addr:           ":8082",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	httpServer.ListenAndServe()

	/*
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
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			cancel()
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
	*/
}
