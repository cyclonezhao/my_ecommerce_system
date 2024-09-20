package client

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

func initXORM(config *DbConfig) *xorm.Engine {
	engine, err := xorm.NewEngine("mysql", config.DataSourceName)
	if err != nil {
		log.Fatalln("数据库连接失败", err)
	}
	return engine
}
