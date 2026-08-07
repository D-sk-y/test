package config

import (
	"os"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	// 数据库配置
	DB struct {
		Driver   string `yaml:"driver"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"db"`
}

func NewConfig() *Config {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}

	return &config
}
