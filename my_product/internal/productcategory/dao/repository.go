package dao

import (
	"my_ecommerce_system/pkg/db"
	"my_product/internal/productcategory/dto"

	my_client "my_ecommerce_system/pkg/client"
	"time"
)

func AddProductCategory(productcategory *dto.ProductCategory) error {
	productcategory.Id = db.GenId()
	productcategory.CreatedAt = time.Now()
	productcategory.UpdatedAt = time.Now()
	_, err := my_client.XORM.Insert(productcategory)
	return err
}

func DeleteProductCategory(id uint64) error {
	sql := `DELETE FROM prod_category WHERE id = ?`
	_, err := my_client.XORM.Exec(sql, id)
	return err
}

func UpdateProductCategory(productcategory *dto.ProductCategory) error {
	productcategory.UpdatedAt = time.Now()
	_, err := my_client.XORM.ID(productcategory.Id).Update(productcategory)
	return err
}

func GetProductCategory(id uint64) (*dto.ProductCategory, error) {
	var productCategory dto.ProductCategory
	_, err := my_client.XORM.ID(id).Get(&productCategory)

	return &productCategory, err
}

func GetProductCategoryList() ([]dto.ProductCategory, error) {
	var productCategoryList []dto.ProductCategory
	err := my_client.XORM.Find(&productCategoryList)
	return productCategoryList, err
}
