package httprouter

import (
	"my_system/internal/user"

	"github.com/gin-gonic/gin"
)

func Route(engine *gin.Engine) {
	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	api := engine.Group("/api")
	{
		api.POST("/user/signUp", user.SignUp)
		api.POST("/user/signIn", user.SignIn)
		api.POST("/user/login", user.SignIn)
		api.GET("/user/signOut", user.SignOut)
	}
}
