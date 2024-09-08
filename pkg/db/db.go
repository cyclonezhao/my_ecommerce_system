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

func Execute(sql string, args ...interface{}){
	stmt, err := db.Prepare(sql)
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
