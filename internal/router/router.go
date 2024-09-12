package router

import (
	"github.com/gin-gonic/gin"
	"my_ecommerce_system/internal/user"
)

func NewRouter() *gin.Engine{
	engine := gin.Default()
	api := engine.Group("/api")
	{
		api.POST("/user/signUp", user.SignUp)
	}
	return engine
}
