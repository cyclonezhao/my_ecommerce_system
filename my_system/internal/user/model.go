package user

import "time"

// 登录请求结构体
type SignUpRequest struct{
	Username string `json:"username"`
	Password string `json:"password"`
}

// 登录响应结构体
type SignUpRespose struct{
	Token string `json:"token"`
}

// 用户实体
type User struct{
	Id uint64 `json:"id"`
	Name string `json:"name"`
	Password string `json:"password"`
	Created_at time.Time `json:"createAt"`
}