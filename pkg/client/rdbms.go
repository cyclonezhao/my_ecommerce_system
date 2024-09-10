package client

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"my_ecommerce_system/pkg/config"
)

var DB *sql.DB

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
	DB = testdb
	DB.SetMaxOpenConns(maxOpenConns) //设置最大打开连接数
	DB.SetMaxIdleConns(maxIdleConns)   //设置最大空闲连接数
	log.Println("数据库初始化成功")
	return nil
}