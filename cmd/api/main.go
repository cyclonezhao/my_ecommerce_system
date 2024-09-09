package main

import (
	"fmt"
	"log"
	"my_ecommerce_system/pkg/config"
	"my_ecommerce_system/pkg/db"
	"my_ecommerce_system/pkg/errorhandler"
	"my_ecommerce_system/pkg/middleware"
	"net/http"

	"github.com/gorilla/mux"

	"my_ecommerce_system/internal/user"
)

func main() {
	startHTTPServer()
}

func startHTTPServer() {
	// 初始化配置
	config.InitConfig()
	// 初始化数据库
	db.InitDB()

	// 配置路由表
	r := mux.NewRouter()
	r.HandleFunc("/user/sayHello", user.SayHello)
	r.Handle("/user/signUp", errorhandler.ErrorToHttpResponse(user.SignUp))
	r.Handle("/user/signIn", errorhandler.ErrorToHttpResponse(user.SignIn))
	r.Handle("/user/login", errorhandler.ErrorToHttpResponse(user.SignIn))
	r.HandleFunc("/hello", hello)

	// 启动http服务
	log.Printf("%s 开始启动！", config.AppConfig.AppName)
	addr := config.AppConfig.Addr
	http.ListenAndServe(addr, middleware.ErrorToHttpHandlingMiddleware(r))
}

func hello(writer http.ResponseWriter, request *http.Request) {
	msg := "Hello, My Ecommerce System 已上线！"
	fmt.Println(msg)
	fmt.Fprintf(writer, msg)
}
