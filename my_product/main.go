package main

import (
	"github.com/gin-gonic/gin"
	"my_ecommerce_system/pkg/web"
	"net/http"
	"time"
)

func main(){
	engine := gin.Default()
	engine.GET("/ping", func(ctx *gin.Context){
		web.Test("===========8081")
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	server := &http.Server{
		Addr:           ":8081",
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
