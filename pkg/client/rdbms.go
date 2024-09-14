package client

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type DbConfig struct {
	DriverName     string `yaml:"driverName"`
	DataSourceName string `yaml:"dataSourceName"`
	MaxOpenConns   int    `yaml:"maxOpenConns"`
	MaxIdleConns   int    `yaml:"maxIdleConns"`
}

func initDB(config *DbConfig) *sql.DB {
	// 数据库类型，后续做成可配置项
	driverName := config.DriverName
	dataSourceName := config.DataSourceName
	maxOpenConns := config.MaxOpenConns
	maxIdleConns := config.MaxIdleConns

	testdb, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatalf("数据库连接失败: %s", err.Error())
	}
	if testdb == nil {
		log.Fatal("数据库打开失败！")
	}

	// 测试连接，理由见sql.Open的方法注释
	err2 := testdb.Ping()
	if err2 != nil {
		log.Fatalf("数据库连接失败: %s", err2.Error())
	}

	// 至此，数据库初始化成功
	DB = testdb
	DB.SetMaxOpenConns(maxOpenConns) //设置最大打开连接数
	DB.SetMaxIdleConns(maxIdleConns) //设置最大空闲连接数
	log.Println("数据库初始化成功")
	return DB
}
