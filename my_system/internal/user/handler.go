package user

import (
	"fmt"
	"my_ecommerce_system/pkg/errorhandler"
	"my_ecommerce_system/pkg/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SayHello(writer http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	msg := fmt.Sprintf("Hi, I am %s", name)
	fmt.Println(msg)
	fmt.Fprintf(writer, msg)
}

func SignUp(ctx *gin.Context) {
	var request SignUpRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		web.ResponseError(ctx, &errorhandler.BusinessError{Message: "请求无效", HttpCode: http.StatusBadRequest})
		return
	}

	tokenString, err := SignUpService(request, new(StdUserRepository))
	if err != nil {
		web.ResponseError(ctx, err)
	} else {
		web.ResponseSuccess(ctx, gin.H{"token": tokenString})
	}
}

func SignIn(ctx *gin.Context) {
	var request SignUpRequest
	err := ctx.ShouldBind(&request)
	if err != nil {
		web.ResponseError(ctx, &errorhandler.BusinessError{Message: "请求无效", HttpCode: http.StatusBadRequest})
		return
	}

	tokenString, err := signIn(request, new(StdUserRepository))
	if err != nil {
		web.ResponseError(ctx, err)
	} else {
		web.ResponseSuccess(ctx, gin.H{"token": tokenString})
	}
}

func SignOut(ctx *gin.Context) {
	ctx.GetHeader("token")
	tokenString := ctx.GetHeader("token")
	userName, err := GetUserNameByToken(tokenString)
	if err != nil {
		web.ResponseError(ctx, err)
	}

	err = signOut(userName)
	if err != nil {
		web.ResponseError(ctx, err)
	} else {
		web.ResponseSuccess(ctx, gin.H{"message": "已登出！"})
	}
}
