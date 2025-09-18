package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	DBConfig DBConfig `yaml:"db_config"`
}

type DBConfig struct {
	Addr string `yaml:"addr"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

var cfg Config

func GetConfig() *Config {
	return &cfg
}

func init() {
	data, err := os.ReadFile("etc/config/config.yaml")
	if err != nil {
		log.Fatalf("读取 config.yaml 失败: %v", err)
	}
	
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("YAML 解析失败: %v", err)
	}
}
