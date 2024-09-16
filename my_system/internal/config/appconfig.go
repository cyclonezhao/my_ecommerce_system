package config

import (
	my_client "my_ecommerce_system/pkg/client"
)

type Config struct {
	DB      my_client.DbConfig    `yaml:"db"`
	Redis   my_client.RedisConfig `yaml:"redis"`
	Gateway struct {
		WriteList []string `yaml:"writeList"`
	} `yaml:"gateway"`
	Jwt struct {
		Expire int `yaml:"expire"`
	} `yaml:"jwt"`
}

var AppConfig Config
