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
		return
	}

	// 根据传来的userName从DB中查用户
	user, err := getUserByName(signInRequest.Username)
	if user == nil{
		http.Error(writer, "用户不存在", http.StatusUnauthorized)
		return
	}else if err != nil{
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		return
	}

	userName := user.Name
	password := user.Password

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(signInRequest.Password)); err != nil{
		http.Error(writer, "用户名或密码错误", http.StatusUnauthorized)
		return
	}

	// 创建 JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userName,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	// 使用密钥签名token
	tokenString, err := token.SignedString([]byte("I_am_a_secretKey"))
	if err != nil{
		http.Error(writer, "创建Token失败", http.StatusInternalServerError)
		return
	}

	// TODO 将token存到Redis


	// 把token返回给前端
	json.NewEncoder(writer).Encode(SignUpRespose{Token: tokenString})
}
