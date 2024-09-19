package service

import (
	"my_product/internal/productcategory/dao"
	"my_product/internal/productcategory/dto"
)

// 涉及外部中间件的调用（如MySQL，Redis等）的方法都放在这里，以便单测时能Mock
type ProductCategoryRepository interface {
	AddProductCategory(productcategory *dto.ProductCategory) error
	DeleteProductCategory(id uint64) error
	UpdateProductCategory(productcategory *dto.ProductCategory) error
	GetProductCategory(id uint64) (*dto.ProductCategory, error)
	GetProductCategoryList() ([]dto.ProductCategory, error)
}

// repository的标准实现
type StdProductCategoryRepository struct{}

func (*StdProductCategoryRepository) AddProductCategory(productcategory *dto.ProductCategory) error {
	return dao.AddProductCategory(productcategory)
}

func (*StdProductCategoryRepository) DeleteProductCategory(id uint64) error {
	return dao.DeleteProductCategory(id)
}

func (*StdProductCategoryRepository) UpdateProductCategory(productcategory *dto.ProductCategory) error {
	return dao.UpdateProductCategory(productcategory)
}

func (*StdProductCategoryRepository) GetProductCategory(id uint64) (*dto.ProductCategory, error) {
	return dao.GetProductCategory(id)
}

func (*StdProductCategoryRepository) GetProductCategoryList() ([]dto.ProductCategory, error) {
	return dao.GetProductCategoryList()
}

// repository的标准实现实例
var StdProductCategoryRepositoryInstance *StdProductCategoryRepository = new(StdProductCategoryRepository)
