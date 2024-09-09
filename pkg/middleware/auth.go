package middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"my_ecommerce_system/pkg/client"
	"my_ecommerce_system/pkg/config"
	"my_ecommerce_system/pkg/constant"
	"net/http"
)

type Claims struct {
	UserName string `json:"username"`
	jwt.RegisteredClaims
}

func AuthenticationMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		// 没想到 Go 竟然没有内置 contains 函数
		path := r.URL.Path
		inWriteList := false
		for _, item := range config.AppConfig.Gateway.WriteList{
			if path == item{
				inWriteList = true
			}
		}

		if(!inWriteList && validateToken(r, w)){
			return
		}

		// token 验证通过或不需要验证
		next.ServeHTTP(w, r)
	})
}

// 校验token，成功返回false，失败返回true（这是IDEA的Extract Method功能自动生成的，懒得改了）
func validateToken(r *http.Request, w http.ResponseWriter) bool {
	tokenString := r.Header.Get("token")
	if tokenString == "" {
		// TODO 这种情况应该跳转回登录页，不过这应该由前端判断到401返回码后执行，或者由后端发301？
		http.Error(w, "token不能为空", http.StatusUnauthorized)
		return true
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(constant.JWT_SECRET_KEY), nil
	})

	// 验证token自身有效性
	if err != nil || !token.Valid {
		msg := "token无效！"
		if err != nil {
			msg += err.Error()
		}
		http.Error(w, msg, http.StatusUnauthorized)
		return true
	}

	// 验证token是否存在
	ctx := context.Background()
	cacheKey := constant.CACHE_USER_TOKEN + ":" + claims.UserName
	storedToken, err := client.RedisClient.Get(ctx, cacheKey).Result()
	if err != nil || storedToken != tokenString {
		msg := "token不存在！"
		if err != nil {
			msg += err.Error()
		}
		http.Error(w, msg, http.StatusUnauthorized)
		return true
	}
	return false
}