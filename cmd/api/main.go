package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"my_ecommerce_system/pkg/client"
	"my_ecommerce_system/pkg/config"
	"net/http"
)

func main() {
	startHTTPServer()
}

func startHTTPServer() {
	// 初始化配置
	config.InitConfig()
	// 初始化数据库
	client.InitDB()
	// 初始化Redis连接
	client.InitRedis()

	// 新代码
	engine := gin.Default()
	engine.GET("/ping", func(ctx *gin.Context){
		ctx.JSON(200, gin.H{"message": "pong"})
	})
	engine.Run()
/*
	// 配置路由表
	r := mux.NewRouter()
	r.HandleFunc("/user/sayHello", user.SayHello)
	r.Handle("/user/signUp", errorhandler.ErrorToHttpResponse(user.SignUp))
	r.Handle("/user/signIn", errorhandler.ErrorToHttpResponse(user.SignIn))
	r.Handle("/user/login", errorhandler.ErrorToHttpResponse(user.SignIn))
	r.Handle("/user/signOut", errorhandler.ErrorToHttpResponse(user.SignOut))
	r.HandleFunc("/hello", hello)

	// 启动http服务
	log.Printf("%s 开始启动！", config.AppConfig.AppName)
	addr := config.AppConfig.Addr

	handler := middleware.AuthenticationMiddleware(
		middleware.ErrorToHttpHandlingMiddleware(r))
	http.ListenAndServe(addr, handler)*/
}

func hello(writer http.ResponseWriter, request *http.Request) {
	msg := "Hello, My Ecommerce System 已上线！"
	fmt.Println(msg)
	fmt.Fprintf(writer, msg)
}
