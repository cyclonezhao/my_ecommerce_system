package productcategory

import (
	"my_ecommerce_system/pkg/errorhandler"
	"my_ecommerce_system/pkg/web"
	"net/http"

	"strconv"

	"encoding/json"

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

	err = AddProductCategoryService(&productcategory, StdProductCategoryRepositoryInstance)
	if err != nil {
		web.ResponseError(ctx, err)
	} else {
		web.ResponseSuccess(ctx, gin.H{"message": "添加成功"})
	}

}

// 删除商品分类
func DeleteProductCategoryHandler(ctx *gin.Context) {
	rawId, exist := ctx.GetQuery("id")
	if !exist {
		web.ResponseError(ctx, &errorhandler.BusinessError{Message: "请求必须包含id参数", HttpCode: http.StatusBadRequest})
		return
	}

	id, err := strconv.ParseUint(rawId, 10, 64)
	if err != nil {
		web.ResponseError(ctx, &errorhandler.BusinessError{Message: "id值必须为整形", HttpCode: http.StatusBadRequest})
		return
	}

	err = DeleteProductCategoryService(id, StdProductCategoryRepositoryInstance)
	if err != nil {
		web.ResponseError(ctx, err)
	} else {
		web.ResponseSuccess(ctx, gin.H{"message": "删除成功"})
	}
}

// 修改商品分类
func UpdateProductCategoryHandler(ctx *gin.Context) {
	var productcategory ProductCategory
	err := ctx.ShouldBind(&productcategory)
	if err != nil {
		web.ResponseError(ctx, &errorhandler.BusinessError{Message: "请求无效", HttpCode: http.StatusBadRequest})
		return
	}

	err = UpdateProductCategoryService(&productcategory, StdProductCategoryRepositoryInstance)
	if err != nil {
		web.ResponseError(ctx, err)
	} else {
		web.ResponseSuccess(ctx, gin.H{"message": "更新成功"})
	}

}

// 查看商品分类
func GetProductCategoryHandler(ctx *gin.Context) {
	rawId, exist := ctx.GetQuery("id")
	if !exist {
		web.ResponseError(ctx, &errorhandler.BusinessError{Message: "请求必须包含id参数", HttpCode: http.StatusBadRequest})
		return
	}

	id, err := strconv.ParseUint(rawId, 10, 64)
	if err != nil {
		web.ResponseError(ctx, &errorhandler.BusinessError{Message: "id值必须为整形", HttpCode: http.StatusBadRequest})
		return
	}

	productCategory, err := GetProductCategoryService(id, StdProductCategoryRepositoryInstance)
	if err != nil {
		web.ResponseError(ctx, err)
	} else {
		b, _ := json.Marshal(productCategory)
		web.ResponseSuccess(ctx, string(b))
	}
}

// 列表商品分类
