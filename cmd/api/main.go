package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"my_ecommerce_system/internal/user"
)

func main() {
	startHTTPServer()
}

func startHTTPServer() {
	r := mux.NewRouter()

	r.HandleFunc("/user/sayHello", user.SayHello)
	r.HandleFunc("/hello", hello)

	addr := ":8080"
	http.ListenAndServe(addr, r)
}

func hello(writer http.ResponseWriter, request *http.Request) {
	msg := "Hello, My Ecommerce System 已上线！"
	fmt.Println(msg)
	fmt.Fprintf(writer, msg)
}
