package productcategory

import (
	"my_ecommerce_system/pkg/errorhandler"
	"my_ecommerce_system/pkg/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 增加商品分类
func AddProductCategoryHandler(ctx *gin.Context) {
	var productcategory ProductCategory
	err := ctx.ShouldBind(&productcategory)
	if err != nil {
		web.ResponseError(ctx, &errorhandler.BusinessError{Message: "请求无效", HttpCode: http.StatusBadRequest})
		return
	}

	err = AddProductCategoryService(&productcategory, new(StdProductCategoryRepository))
	if err != nil {
		web.ResponseError(ctx, err)
	} else {
		web.ResponseSuccess(ctx, gin.H{"message": "添加成功"})
	}

}

// 删除商品分类
// 修改商品分类
// 查看商品分类
// 列表商品分类
