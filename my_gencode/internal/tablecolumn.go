package internal

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Table struct {
	Name    string
	Comment string
}

type TableColumn struct {
	ColumnName    string
	DataType      string
	ColumnKey     string
	ColumnComment string
}

func GetTableInfo(tableName string) (table *Table, tableColumns []*TableColumn) {
	db := initDB()
	table = getTable(db, tableName)
	tableColumns = getTableColumns(db, tableName)
	return
}

func getTable(db *sql.DB, tableName string) *Table {
	sqlStr := `SELECT TABLE_NAME, TABLE_COMMENT 
		FROM information_schema.TABLES WHERE TABLE_NAME = ? `

	var table Table
	err := executeQuery(db, sqlStr, func(rows *sql.Rows) error {
		for rows.Next() {
			err := rows.Scan(&table.Name, &table.Comment)
			if err != nil {
				return err
			}
			break
		}
		return nil
	}, tableName)

	if err != nil {
		log.Fatalln("解析表信息失败！", err)
	}

	return &table
}

func getTableColumns(db *sql.DB, tableName string) []*TableColumn {
	sqlStr := `SELECT COLUMN_NAME, DATA_TYPE, COLUMN_KEY, COLUMN_COMMENT 
		FROM information_schema.COLUMNS WHERE TABLE_NAME = ? `

	var tableColumns []*TableColumn
	err := executeQuery(db, sqlStr, func(rows *sql.Rows) error {
		for rows.Next() {
			var tableColumn TableColumn
			err := rows.Scan(&tableColumn.ColumnName, &tableColumn.DataType, &tableColumn.ColumnKey, &tableColumn.ColumnComment)
			if err != nil {
				return err
			}

			tableColumns = append(tableColumns, &tableColumn)
		}
		return nil
	}, tableName)

	if err != nil {
		log.Fatalln("解析表结构失败！", err)
	}

	return tableColumns
}

func initDB() *sql.DB {
	// 数据库类型，后续做成可配置项
	driverName := "mysql"
	dataSourceName := "root:root@/my_ecommerce_system"

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
	log.Println("数据库初始化成功")
	return testdb
}

func executeQuery(db *sql.DB, query string, rowsHandler func(rows *sql.Rows) error, args ...interface{}) error {
	rows, err := db.Query(query, args...)
	if err != nil {
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
