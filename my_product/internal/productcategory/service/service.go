package service

import (
	"fmt"
	"my_product/internal/productcategory/dto"
)

func AddProductCategory(productcategory *dto.ProductCategory, repository ProductCategoryRepository) error {
	return repository.AddProductCategory(productcategory)
}

func DeleteProductCategory(id uint64, repository ProductCategoryRepository) error {
	return repository.DeleteProductCategory(id)
}

func UpdateProductCategory(productcategory *dto.ProductCategory, repository ProductCategoryRepository) error {
	if productcategory.Id == 0 {
		return fmt.Errorf("更新操作，id不能为空！")
	}
	return repository.UpdateProductCategory(productcategory)
}

func GetProductCategory(id uint64, repository ProductCategoryRepository) (*dto.ProductCategory, error) {
	return repository.GetProductCategory(id)
}

func GetProductCategoryList(repository ProductCategoryRepository) ([]dto.ProductCategory, error) {
	return repository.GetProductCategoryList()
}
