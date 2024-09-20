package handler

import (
	"github.com/gin-gonic/gin"
)

func Route(api *gin.RouterGroup) {
	productCategory := api.Group("productCategory")

	productCategory.POST("/AddProductCategory", AddProductCategoryHandler)
	productCategory.DELETE("/DeleteProductCategory", DeleteProductCategoryHandler)
	productCategory.PUT("/UpdateProductCategory", UpdateProductCategoryHandler)
	productCategory.GET("/GetProductCategory", GetProductCategoryHandler)
	productCategory.GET("/GetProductCategoryList", GetProductCategoryHandlerList)
}
