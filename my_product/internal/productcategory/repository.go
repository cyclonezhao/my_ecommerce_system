package productcategory

import (
	"my_ecommerce_system/pkg/db"

	"time"
)

func AddProductCategory(productcategory *ProductCategory) error {
	sql := `INSERT INTO prod_category (id, name, created_at) VALUES (?, ?, ?)`
	return db.Execute(sql, db.GenId(), productcategory.Name, time.Now())
}
