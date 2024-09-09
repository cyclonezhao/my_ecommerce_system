package user

import (
	"encoding/json"
	"fmt"
	"my_ecommerce_system/pkg/errorhandler"
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

func SignUp(writer http.ResponseWriter, request *http.Request) error{
	var signInRequest SignUpRequest
	if err := json.NewDecoder(request.Body).Decode(&signInRequest); err != nil {
		return &errorhandler.BusinessError{Message:"请求无效", HttpCode:http.StatusBadRequest}
	}

	err := signUp(signInRequest)

	if err != nil{
		// 创建Token
		tokenString, err := genToken(signInRequest.Username)
		if err != nil{
			return err
		}

		// 把token返回给前端
		json.NewEncoder(writer).Encode(SignUpRespose{Token: tokenString})
	}
	return err
}

func SignIn(writer http.ResponseWriter, request *http.Request) error {
	var signInRequest SignUpRequest
	if err := json.NewDecoder(request.Body).Decode(&signInRequest); err != nil {
		return &errorhandler.BusinessError{Message:"请求无效", HttpCode:http.StatusBadRequest}
	}

	// 根据传来的userName从DB中查用户
	user, err := getUserByName(signInRequest.Username)
	if user == nil{
		return &errorhandler.BusinessError{Message:"用户不存在", HttpCode:http.StatusUnauthorized}
	}else if err != nil{
		return &errorhandler.BusinessError{Message:err.Error(), HttpCode:http.StatusUnauthorized}
	}

	userName := user.Name
	password := user.Password

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(signInRequest.Password)); err != nil{
		return &errorhandler.BusinessError{Message:"用户名或密码错误", HttpCode:http.StatusUnauthorized}
	}

	// 创建Token
	tokenString, err := genToken(userName)
	if err != nil{
		return err
	}

	// 把token返回给前端
	json.NewEncoder(writer).Encode(SignUpRespose{Token: tokenString})
	return nil
}

func genToken(userName string) (string, error){
	// 创建 JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userName,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	// 使用密钥签名token
	tokenString, err := token.SignedString([]byte("I_am_a_secretKey"))
	if err != nil{
		return "", &errorhandler.BusinessError{Message:"创建Token失败", HttpCode:http.StatusUnauthorized}
	}

	// TODO 将token存到Redis

	return tokenString, nil
}
