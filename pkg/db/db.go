package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"my_ecommerce_system/pkg/config"
)

// 全局私有单例
var idGenerator = newIDGenerator()
var db *sql.DB

func GenId() uint64{
	return idGenerator.generateID()
}

func InitDB() error {
	// 数据库类型，后续做成可配置项
	driverName := config.AppConfig.DB.DriverName
	dataSourceName := config.AppConfig.DB.DataSourceName
	maxOpenConns := config.AppConfig.DB.MaxOpenConns
	maxIdleConns := config.AppConfig.DB.MaxIdleConns

	testdb, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatal("数据库连接失败: %s", err.Error())
	}
	if testdb == nil {
		log.Fatal("数据库打开失败！")
	}

	// 测试连接，理由见sql.Open的方法注释
	err2 := testdb.Ping()
	if err2 != nil {
		log.Fatal("数据库连接失败: %s", err2.Error())
	}

	// 至此，数据库初始化成功
	db = testdb
	db.SetMaxOpenConns(maxOpenConns) //设置最大打开连接数
	db.SetMaxIdleConns(maxIdleConns)   //设置最大空闲连接数
	log.Println("数据库初始化成功")
	return nil
}

func Execute(query string, args ...interface{}) error{
	stmt, err := db.Prepare(query)
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
	rows, err := db.Query(query, args...)
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
	err := db.QueryRow(query, args...).Scan(&exists)
	if err != nil{
		return false, err
	}
	return exists, nil
}
