package db

import (
	"database/sql"
	"fmt"
	"my_ecommerce_system/pkg/client"
)

// 全局私有单例
var idGenerator = newIDGenerator()

func GenId() uint64{
	return idGenerator.generateID()
}

func Execute(query string, args ...interface{}) error{
	stmt, err := client.DB.Prepare(query)
	if err != nil{
		return err
	}

	res, err := stmt.Exec(args...)
	if err != nil{
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	fmt.Println("rowsAffected", rowsAffected)
	return nil
}

func ExecuteQuery(query string, rowsHandler func(rows *sql.Rows) error, args ...interface{}) error{
	rows, err := client.DB.Query(query, args...)
	if err != nil{
		return err
	}
	defer rows.Close()

	err = rowsHandler(rows)
	if err != nil {
		return err
	}

	// 在遍历 rows 后检查 rows.Err() 以确保没有发生错误
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

func Exists(query string, args ...interface{}) (bool, error){
	var exists bool
	query = fmt.Sprintf("SELECT EXISTS (%s)", query)
	err := client.DB.QueryRow(query, args...).Scan(&exists)
	if err != nil{
		return false, err
	}
	return exists, nil
}
