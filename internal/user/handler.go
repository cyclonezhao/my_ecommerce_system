package user

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"my_ecommerce_system/pkg/errorhandler"
	"net/http"
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
	if err != nil{
		// TODO 请求无效
	}

	tokenString, err := SignUpService(request, new(StdUserRepository))

	ctx.JSON(200, gin.H{"message": tokenString})
}

func SignIn(writer http.ResponseWriter, request *http.Request) error {
	var signInRequest SignUpRequest
	if err := json.NewDecoder(request.Body).Decode(&signInRequest); err != nil {
		return &errorhandler.BusinessError{Message:"请求无效", HttpCode:http.StatusBadRequest}
	}

	tokenString, err := signIn(signInRequest, new(StdUserRepository))
	if err != nil{
		return err
	}

	// 把token返回给前端
	json.NewEncoder(writer).Encode(SignUpRespose{Token: tokenString})
	return nil
}

func SignOut(writer http.ResponseWriter, request *http.Request) error {
	userName := request.FormValue("userName")
	err := signOut(userName)
	if err == nil{
		writer.Write([]byte("已登出！"))
	}
	return err
}
