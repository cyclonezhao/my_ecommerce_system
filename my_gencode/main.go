package main

import (
	"bytes"
	"io/fs"
	"log"
	"my_gencode/internal"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

type DataToTemplate struct {
	EntityName  string
	Columns     []*internal.TableColumn
	Table       *internal.Table
	ProjectName string
}

func main() {
	// 表名
	tableName := "prod_category"
	// 实体名
	entityName := "ProductCategory2"
	// 模块名
	projectName := "my_product"
	// 生成的代码路径前缀
	targetPathPrefix := "C:/Users/light/IdeaProjects/my_ecommerce_system/my_product/internal/"

	// 从information_schema取表信息
	table, tableColumns := internal.GetTableInfo(tableName)

	// 总体数据
	var dataToTemplate = &DataToTemplate{
		entityName, tableColumns, table, projectName,
	}

	// 渲染工具函数集
	funcMap := template.FuncMap{
		"SnakeToUpperCamelCase": SnakeToUpperCamelCase,
		"SnakeToLowerCamelCase": SnakeToLowerCamelCase,
		"GetGoTypeByDBType":     GetGoTypeByDBType,
		"ToLower":               strings.ToLower,
	}

	templatePathPrefix := "./resources/templates/"
	var templatePath string

	templatePath = filepath.Join(templatePathPrefix, "dto", "model")
	render(templatePath, dataToTemplate, targetPathPrefix, funcMap, "dto", "model.go")

	templatePath = filepath.Join(templatePathPrefix, "dao", "repository")
	render(templatePath, dataToTemplate, targetPathPrefix, funcMap, "dao", "repository.go")

	templatePath = filepath.Join(templatePathPrefix, "handler", "handler")
	render(templatePath, dataToTemplate, targetPathPrefix, funcMap, "handler", "handler.go")

	templatePath = filepath.Join(templatePathPrefix, "handler", "router")
	render(templatePath, dataToTemplate, targetPathPrefix, funcMap, "handler", "router.go")

	templatePath = filepath.Join(templatePathPrefix, "service", "repositoryinterface")
	render(templatePath, dataToTemplate, targetPathPrefix, funcMap, "service", "repositoryinterface.go")

	templatePath = filepath.Join(templatePathPrefix, "service", "service")
	render(templatePath, dataToTemplate, targetPathPrefix, funcMap, "service", "service.go")
}

func render(templatePath string, dataToTemplate *DataToTemplate, targetPathPrefix string, funcMap template.FuncMap, targetPathEnd ...string) {
	// 读取模板文件内容
	content, err := os.ReadFile(templatePath)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	// 渲染代码文件
	tpl := template.Must(template.New("model").Funcs(funcMap).Parse(string(content)))

	var buf bytes.Buffer
	tpl.Execute(&buf, dataToTemplate)

	// 构建输出目录以及输出代码文件
	targetPath := filepath.Join(targetPathPrefix, strings.ToLower(dataToTemplate.EntityName), filepath.Join(targetPathEnd...))

	// 确保目录存在，如果不存在则创建
	err = os.MkdirAll(filepath.Dir(targetPath), fs.ModePerm)
	if err != nil {
		log.Printf("创建目录失败: %v\n", err)
		return
	}
	// 将渲染结果写入文件
	err = os.WriteFile(targetPath, buf.Bytes(), fs.ModePerm)
	if err != nil {
		log.Printf("写入文件失败: %v\n", err)
		return
	}
}

func SnakeToUpperCamelCase(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Title(s)
	return strings.Replace(s, " ", "", -1)
}

func SnakeToLowerCamelCase(s string) string {
	s = SnakeToUpperCamelCase(s)
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}

func GetGoTypeByDBType(dbType string) string {
	return DBTypeToStructType[dbType]
}

// 数据库字段类型，到Go变量类型的映射
var DBTypeToStructType = map[string]string{
	"int":        "int32",
	"tinyint":    "int8",
	"smallint":   "int",
	"mediumint":  "uint64",
	"bigint":     "uint64",
	"bit":        "int",
	"bool":       "bool",
	"enum":       "string",
	"set":        "string",
	"varchar":    "string",
	"char":       "string",
	"tinytext":   "string",
	"mediumtext": "string",
	"text":       "string",
	"longtext":   "string",
	"blob":       "string",
	"tinyblob":   "string",
	"mediumblob": "string",
	"longblob":   "string",
	"date":       "time.Time",
	"datetime":   "time.Time",
	"timestamp":  "time.Time",
	"time":       "time.Time",
	"float":      "float64",
	"double":     "float64",
}
