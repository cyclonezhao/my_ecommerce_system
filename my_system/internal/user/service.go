package user

import (
	"context"
	"fmt"
	"log"
	"my_ecommerce_system/pkg/client"
	"my_ecommerce_system/pkg/constant"
	"my_ecommerce_system/pkg/errorhandler"
	"my_ecommerce_system/pkg/middleware"
	"my_system/internal/config"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// 涉及外部中间件的调用（如MySQL，Redis等）的方法都放在这里，以便单测时能Mock
type UserRepository interface {
	ExistsUserName(userName string) (bool, error)
	AddNewUser(user *User) error
	GetTokenExpirationTime() time.Duration
	SetTokenIntoRedis(cacheKey string, tokenString string, expirationTime time.Duration) error
}

type StdUserRepository struct{}

func (r *StdUserRepository) AddNewUser(newUser *User) error {
	return AddNewUser(newUser)
}

func (r *StdUserRepository) ExistsUserName(name string) (bool, error) {
	return ExistsUserName(name)
}

func (r *StdUserRepository) GetTokenExpirationTime() time.Duration {
	return time.Duration(config.AppConfig.Jwt.Expire) * time.Second
}

func (r *StdUserRepository) SetTokenIntoRedis(cacheKey string, tokenString string, expirationTime time.Duration) error {
	ctx := context.Background()
	return client.RedisClient.Set(ctx, cacheKey, tokenString, expirationTime).Err()
}

func SignUpService(request SignUpRequest, repository UserRepository) (string, error) {
	// 用户名，密码
	userName := request.Username
	password := request.Password

	// 检查用户是否存在
	exists, err := repository.ExistsUserName(userName)
	if err != nil {
		return "", err
	} else if exists {
		return "", &errorhandler.BusinessError{
			Message: fmt.Sprintf("用户名[%s]已存在！", userName),
		}
	}

	// 密码哈希化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("密码明文哈希化失败")
	}

	// 创建用户
	newUser := &User{
		Name:     userName,
		Password: string(hashedPassword),
	}
	repository.AddNewUser(newUser)

	// 创建Token
	tokenString, err := genToken(userName, repository)
	return tokenString, err
}

func signIn(signInRequest SignUpRequest, repository UserRepository) (string, error) {
	// 根据传来的userName从DB中查用户
	user, err := getUserByName(signInRequest.Username)
	if user == nil {
		return "", &errorhandler.BusinessError{Message: "用户不存在", HttpCode: http.StatusUnauthorized}
	} else if err != nil {
		return "", &errorhandler.BusinessError{Message: err.Error(), HttpCode: http.StatusUnauthorized}
	}

	userName := user.Name
	password := user.Password

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(signInRequest.Password)); err != nil {
		return "", &errorhandler.BusinessError{Message: "用户名或密码错误", HttpCode: http.StatusUnauthorized}
	}

	// 创建Token
	tokenString, err := genToken(userName, repository)
	return tokenString, err
}

func genToken(userName string, repository UserRepository) (string, error) {
	// 创建 JWT Token
	expirationTime := repository.GetTokenExpirationTime()
	claims := &middleware.Claims{
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名token
	tokenString, err := token.SignedString([]byte(constant.JWT_SECRET_KEY))
	if err != nil {
		return "", &errorhandler.BusinessError{Message: "创建Token失败"}
	}

	// 将token存到Redis
	cacheKey := constant.CACHE_USER_TOKEN + ":" + userName
	err = repository.SetTokenIntoRedis(cacheKey, tokenString, expirationTime)
	if err != nil {
		log.Print(err)
		return "", &errorhandler.BusinessError{Message: "缓存Token失败"}
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
