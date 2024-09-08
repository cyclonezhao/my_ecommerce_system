package user

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	//"net/http"
)

func signUp(request SignUpRequest) error{
	// 用户名，密码
	userName := request.Username
	password := request.Password

	// 检查用户是否存在

	// 密码哈希化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return fmt.Errorf("密码明文哈希化失败")
	}

	// 创建用户
	newUser := &User{
		Name:userName,
		Password: string(hashedPassword),
	}
	addNewUser(newUser)
	return nil
}
