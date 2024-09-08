package main

import (
	"fmt"
	"my_ecommerce_system/pkg/db"
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
	r.HandleFunc("/user/signUp", user.SignUp)
	r.HandleFunc("/user/signIn", user.SignIn)
	r.HandleFunc("/user/login", user.SignIn)
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
