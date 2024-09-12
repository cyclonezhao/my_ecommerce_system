package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"my_ecommerce_system/pkg/errorhandler"
	"net/http"
)

// 返回成功响应，自定义响应体
func ResponseSuccess(ctx *gin.Context, data interface{}){
	if data == nil{
		data = gin.H{}
	}
	ctx.JSON(http.StatusOK, data)
}

// 返回失败响应
func ResponseError(ctx *gin.Context, data interface{}){
	httpStatus := http.StatusInternalServerError
	code := http.StatusInternalServerError
	msg := fmt.Sprintf("Internal Server Error: %v", data)
	if businessError, ok := data.(*errorhandler.BusinessError); ok {
		if businessError.HttpCode != 0{
			httpStatus = businessError.HttpCode
		}

		if businessError.Code != 0{
			code = businessError.Code
		}else{
			code = httpStatus
		}

		msg = businessError.Error()
	}else if err, ok := data.(error); ok{
		msg = err.Error()
	}


	responseData := gin.H{"code": code, "msg": msg}
	ctx.JSON(httpStatus, responseData)
}