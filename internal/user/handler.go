package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func SayHello(writer http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	msg := fmt.Sprintf("Hi, I am %s", name)
	fmt.Println(msg)
	fmt.Fprintf(writer, msg)
}

func SignUp(writer http.ResponseWriter, request *http.Request){
	var signInRequest SignUpRequest
	if err := json.NewDecoder(request.Body).Decode(&signInRequest); err != nil {
		http.Error(writer, "请求无效", http.StatusBadRequest)
	}

	signUp(signInRequest)
}

func SignIn(writer http.ResponseWriter, request *http.Request){
	var signInRequest SignUpRequest
	if err := json.NewDecoder(request.Body).Decode(&signInRequest); err != nil {
		http.Error(writer, "请求无效", http.StatusBadRequest)
	}

	// 临时测试，正确应该从DB中查
	testUsername := signInRequest.Username
	// 明文：123456
	testPassword := "$2a$10$bWpBrhUJKGmcPNc1UB3Fxus0ZNOQtCBmgwWcXMYDyeyC.1H/Ef29G"

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(testPassword), []byte(signInRequest.Password)); err != nil{
		http.Error(writer, "用户名或密码错误", http.StatusUnauthorized)
	}

	// 创建 JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": testUsername,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	// 使用密钥签名token
	tokenString, err := token.SignedString([]byte("I_am_a_secretKey"))
	if err != nil{
		http.Error(writer, "创建Token失败", http.StatusInternalServerError)
	}

	// TODO 将token存到Redis


	// 把token返回给前端
	json.NewEncoder(writer).Encode(SignUpRespose{Token: tokenString})
}
