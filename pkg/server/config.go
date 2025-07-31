package server

import (
	"fmt"

	"github.com/spf13/viper"
)

// LogstashConfig 配置文件结构
type LogstashConfig struct {
	Endpoints []string `mapstructure:"endpoints"` // Logstash API 地址列表
	Web struct {
		ListenAddress string `mapstructure:"listen_address"` // Web 监听地址
	} `mapstructure:"web"`
}

// LoadConfig 从文件加载配置
func LoadConfig(filename string) (*LogstashConfig, error) {
	v := viper.New()

	// 设置配置文件路径
	v.SetConfigFile(filename)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config LogstashConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
}