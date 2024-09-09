package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// 全局私有单例
var idGenerator = newIDGenerator()
var db *sql.DB

func GenId() uint64{
	return idGenerator.generateID()
}

func InitDB() error {
	// 数据库类型，后续做成可配置项
	driverName := "mysql"
	dataSourceName := "root:root@/my_ecommerce_system"
	maxOpenConns := 2000
	maxIdleConns := 1000

	testdb, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return fmt.Errorf("数据库连接失败: %s", err.Error())
	}
	if testdb == nil {
		fmt.Errorf("数据库打开失败！")
	}

	// 测试连接，理由见sql.Open的方法注释
	err2 := testdb.Ping()
	if err2 != nil {
		return fmt.Errorf("数据库连接失败: %s", err2.Error())
	}

	// 至此，数据库初始化成功
	db = testdb
	db.SetMaxOpenConns(maxOpenConns) //设置最大打开连接数
	db.SetMaxIdleConns(maxIdleConns)   //设置最大空闲连接数
	fmt.Println("数据库初始化成功")
	return nil
}

func Execute(query string, args ...interface{}){
	stmt, err := db.Prepare(query)
	if err != nil{
		panic(err.Error())
	}

	res, err := stmt.Exec(args...)
	if err != nil{
		panic(err.Error())
	}

	rowsAffected, _ := res.RowsAffected()
	fmt.Println("rowsAffected", rowsAffected)
}

func ExecuteQuery(query string, rowsHandler func(rows *sql.Rows) error, args ...interface{}){
	rows, err := db.Query(query, args...)
	if err != nil{
		panic(err.Error())
	}
	defer rows.Close()

	err = rowsHandler(rows)
	if err != nil {
		panic(err.Error())
	}

	// 在遍历 rows 后检查 rows.Err() 以确保没有发生错误
	if err := rows.Err(); err != nil {
		panic(err.Error())
	}
}

func Exists(query string, args ...interface{}) (bool, error){
	var exists bool
	query = fmt.Sprintf("SELECT EXISTS (%s)", query)
	err := db.QueryRow(query, args...).Scan(&exists)
	if err != nil{
		return false, err
	}
	return exists, nil
}
