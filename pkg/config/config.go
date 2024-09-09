package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Config struct {
	AppName string `yaml:"appName"`
	Addr    string    `yaml:"addr"`
	DB      struct {
		DriverName     string `yaml:"driverName"`
		DataSourceName     string    `yaml:"dataSourceName"`
		MaxOpenConns     int `yaml:"maxOpenConns"`
		MaxIdleConns int `yaml:"maxIdleConns"`
	} `yaml:"db"`
	Redis struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"redis"`
}

// 全局变量
var AppConfig *Config

func InitConfig(){
	configPath := "./resources/config.yaml"

	if err := loadConfig(configPath); err != nil {
		log.Fatalf("加载配置文件失败：%v", err)
	}
	log.Println("加载配置成功！")
}

func loadConfig(configPath string) error{
	data, err := ioutil.ReadFile(configPath)
	if err != nil{
		return err
	}

	err = yaml.Unmarshal(data, &AppConfig)
	if err != nil{
		return err
	}
	return nil
}