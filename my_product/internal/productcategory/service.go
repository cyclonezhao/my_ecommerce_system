package productcategory

import "fmt"

// 涉及外部中间件的调用（如MySQL，Redis等）的方法都放在这里，以便单测时能Mock
type ProductCategoryRepository interface {
	AddProductCategory(productcategory *ProductCategory) error
	DeleteProductCategory(id uint64) error
	UpdateProductCategory(productcategory *ProductCategory) error
	GetProductCategory(id uint64) (*ProductCategory, error)
	GetProductCategoryList() ([]ProductCategory, error)
}

// repository的标准实现
type StdProductCategoryRepository struct{}

func (*StdProductCategoryRepository) AddProductCategory(productcategory *ProductCategory) error {
	return AddProductCategory(productcategory)
}

func (*StdProductCategoryRepository) DeleteProductCategory(id uint64) error {
	return DeleteProductCategory(id)
}

func (*StdProductCategoryRepository) UpdateProductCategory(productcategory *ProductCategory) error {
	return UpdateProductCategory(productcategory)
}

func (*StdProductCategoryRepository) GetProductCategory(id uint64) (*ProductCategory, error) {
	return GetProductCategory(id)
}

func (*StdProductCategoryRepository) GetProductCategoryList() ([]ProductCategory, error) {
	return GetProductCategoryList()
}

// repository的标准实现实例
var StdProductCategoryRepositoryInstance *StdProductCategoryRepository = new(StdProductCategoryRepository)

func AddProductCategoryService(productcategory *ProductCategory, repository ProductCategoryRepository) error {
	return repository.AddProductCategory(productcategory)
}

func DeleteProductCategoryService(id uint64, repository ProductCategoryRepository) error {
	return repository.DeleteProductCategory(id)
}

func UpdateProductCategoryService(productcategory *ProductCategory, repository ProductCategoryRepository) error {
	if productcategory.Id == 0 {
		return fmt.Errorf("更新操作，id不能为空！")
	}
	return repository.UpdateProductCategory(productcategory)
}

func GetProductCategoryService(id uint64, repository ProductCategoryRepository) (*ProductCategory, error) {
	return repository.GetProductCategory(id)
}

func GetProductCategoryServiceList(repository ProductCategoryRepository) ([]ProductCategory, error) {
	return repository.GetProductCategoryList()
}
