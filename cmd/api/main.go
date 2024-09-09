package main

import (
	"fmt"
	"my_ecommerce_system/pkg/db"
	"my_ecommerce_system/pkg/errorhandler"
	"net/http"

	"github.com/gorilla/mux"

	"my_ecommerce_system/internal/user"
)

func main() {
	startHTTPServer()
}

func startHTTPServer() {
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
	addr := ":8080"
	http.ListenAndServe(addr, r)
}

func hello(writer http.ResponseWriter, request *http.Request) {
	msg := "Hello, My Ecommerce System 已上线！"
	fmt.Println(msg)
	fmt.Fprintf(writer, msg)
}
