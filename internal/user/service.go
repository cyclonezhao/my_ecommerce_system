package user

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"my_ecommerce_system/pkg/client"
	"my_ecommerce_system/pkg/config"
	"my_ecommerce_system/pkg/constant"
	"my_ecommerce_system/pkg/errorhandler"
	"my_ecommerce_system/pkg/middleware"
	"net/http"
	"time"
)

func signUp(request SignUpRequest) (string, error) {
	// 用户名，密码
	userName := request.Username
	password := request.Password

	// 检查用户是否存在
	exists, err := existsUserName(userName)
	if err != nil{
		return "", err
	}else if exists{
		return "", &errorhandler.BusinessError{
			Message:fmt.Sprintf("用户名[%s]已存在！", userName),
		}
	}

	// 密码哈希化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return "", fmt.Errorf("密码明文哈希化失败")
	}

	// 创建用户
	newUser := &User{
		Name:userName,
		Password: string(hashedPassword),
	}
	addNewUser(newUser)

	// 创建Token
	tokenString, err := genToken(userName)
	return tokenString, err
}

func signIn(signInRequest SignUpRequest) (string, error){
	// 根据传来的userName从DB中查用户
	user, err := getUserByName(signInRequest.Username)
	if user == nil{
		return "", &errorhandler.BusinessError{Message:"用户不存在", HttpCode:http.StatusUnauthorized}
	}else if err != nil{
		return "", &errorhandler.BusinessError{Message:err.Error(), HttpCode:http.StatusUnauthorized}
	}

	userName := user.Name
	password := user.Password

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(signInRequest.Password)); err != nil{
		return "", &errorhandler.BusinessError{Message:"用户名或密码错误", HttpCode:http.StatusUnauthorized}
	}

	// 创建Token
	tokenString, err := genToken(userName)
	return tokenString, err
}

func genToken(userName string) (string, error){
	// 创建 JWT Token
	expirationTime := time.Duration(config.AppConfig.Jwt.Expire) * time.Second
	claims := &middleware.Claims{
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名token
	tokenString, err := token.SignedString([]byte(constant.JWT_SECRET_KEY))
	if err != nil{
		return "", &errorhandler.BusinessError{Message:"创建Token失败"}
	}

	// 将token存到Redis
	ctx := context.Background()
	cacheKey := constant.CACHE_USER_TOKEN + ":" + userName
	err = client.RedisClient.Set(ctx, cacheKey, tokenString, expirationTime).Err()
	if err != nil {
		log.Print(err)
		return "", &errorhandler.BusinessError{Message:"缓存Token失败"}
	}

	return tokenString, nil
}

func signOut(userName string) error {
	ctx := context.Background()
	cacheKey := constant.CACHE_USER_TOKEN + ":" + userName

	// 返回的第一个值表示成功删除的key的数量。
	// 如果没有删除成功一个key，就表明cacheKey已经不存在了（可能自动过期了）
	// 所以登出场景不需要关心这个信息
	_, err := client.RedisClient.Del(ctx, cacheKey).Result()
	return err
}