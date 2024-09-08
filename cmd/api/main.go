package main

import (
	"fmt"
	"net/http"
)

func main() {
	startHTTPServer()
}

func startHTTPServer() {
	addr := ":8080"
	http.HandleFunc("/hello", hello)
	http.ListenAndServe(addr, nil)
}

func hello(writer http.ResponseWriter, request *http.Request) {
	msg := "Hello, My Ecommerce System 已上线！"
	fmt.Println(msg)
	fmt.Fprintf(writer, msg)
}
