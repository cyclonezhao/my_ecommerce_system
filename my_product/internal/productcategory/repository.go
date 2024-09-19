package productcategory

import (
	"my_ecommerce_system/pkg/db"

	"database/sql"
	"my_ecommerce_system/pkg/constant"
	"time"
)

func AddProductCategory(productcategory *ProductCategory) error {
	sql := `INSERT INTO prod_category (id, name, created_at, updated_at) VALUES (?, ?, ?, ?)`
	return db.Execute(sql, db.GenId(), productcategory.Name, time.Now(), time.Now())
}

func DeleteProductCategory(id uint64) error {
	sql := `DELETE FROM prod_category WHERE id = ?`
	return db.Execute(sql, id)
}

func UpdateProductCategory(productcategory *ProductCategory) error {
	sql := `UPDATE prod_category SET name = ?, updated_at = ? WHERE id = ?) VALUES (?, ?, ?)`
	return db.Execute(sql, productcategory.Name, time.Now(), productcategory.Id)
}

func GetProductCategory(id uint64) (*ProductCategory, error) {
	sqlStr := `SELECT id, name, created_at, updated_at FROM prod_category WHERE id = ?`
	productCategoryList, err := queryProductCategory(&sqlStr, id)

	if err != nil {
		return nil, err
	}

	return &productCategoryList[0], nil
}

func GetProductCategoryList() ([]ProductCategory, error) {
	sqlStr := `SELECT id, name, created_at, updated_at FROM prod_category`
	return queryProductCategory(&sqlStr)
}

func queryProductCategory(sqlStr *string, args ...interface{}) ([]ProductCategory, error) {
	var productCategoryList []ProductCategory

	// 貌似当前的数据库驱动无法自动将 []uint8 转换为 time.Time
	// 故手动进行转换一下
	// 这是这样一来就破坏了代码统一性
	var createdAt sql.NullString
	var updatedAt sql.NullString

	err := db.ExecuteQuery(*sqlStr, func(rows *sql.Rows) error {
		for rows.Next() {
			var productCategory ProductCategory
			err := rows.Scan(&productCategory.Id, &productCategory.Name, &createdAt, &updatedAt)
			if err != nil {
				return err
			}

			if createdAt.Valid {
				// 很奇怪的时间日期格式化模板：constant.DATE_TIME_FORMAT
				productCategory.Created_at, err = time.Parse(constant.DATE_TIME_FORMAT, createdAt.String)
				if err != nil {
					return err
				}
			}

			if updatedAt.Valid {
				// 很奇怪的时间日期格式化模板：constant.DATE_TIME_FORMAT
				productCategory.Updated_at, err = time.Parse(constant.DATE_TIME_FORMAT, updatedAt.String)
				if err != nil {
					return err
				}
			}

			productCategoryList = append(productCategoryList, productCategory)
		}
		return nil
	}, args...)

	if err != nil {
		return nil, err
	}

	return productCategoryList, nil
}
