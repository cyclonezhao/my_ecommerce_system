package productcategory

// 涉及外部中间件的调用（如MySQL，Redis等）的方法都放在这里，以便单测时能Mock
type ProductCategoryRepository interface {
	AddProductCategory(productcategory *ProductCategory) error
	DeleteProductCategory(id uint64) error
}

// repository的标准实现
type StdProductCategoryRepository struct{}

func (*StdProductCategoryRepository) AddProductCategory(productcategory *ProductCategory) error {
	return AddProductCategory(productcategory)
}

func (*StdProductCategoryRepository) DeleteProductCategory(id uint64) error {
	return DeleteProductCategory(id)
}

// repository的标准实现实例
var StdProductCategoryRepositoryInstance *StdProductCategoryRepository = new(StdProductCategoryRepository)

func AddProductCategoryService(productcategory *ProductCategory, repository ProductCategoryRepository) error {
	return repository.AddProductCategory(productcategory)
}

func DeleteProductCategoryService(id uint64, repository ProductCategoryRepository) error {
	return repository.DeleteProductCategory(id)
}
