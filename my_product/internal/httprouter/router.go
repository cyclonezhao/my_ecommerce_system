package httprouter

import (
	productcategory "my_product/internal/productcategory/handler"

	"github.com/gin-gonic/gin"
)

func Route(engine *gin.Engine) {
	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	api := engine.Group("/api")
	productcategory.Route(api)
}
