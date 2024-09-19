package productcategory

// 涉及外部中间件的调用（如MySQL，Redis等）的方法都放在这里，以便单测时能Mock
type ProductCategoryRepository interface {
	AddProductCategory(productcategory *ProductCategory) error
}

// repository的标准实现
type StdProductCategoryRepository struct{}

func (*StdProductCategoryRepository) AddProductCategory(productcategory *ProductCategory) error {
	return AddProductCategory(productcategory)
}

func AddProductCategoryService(productcategory *ProductCategory, repository ProductCategoryRepository) error {
	return repository.AddProductCategory(productcategory)
}
