package cmd

import (
	"fmt"
	"os"
	"strings"

	"go-logstash-exporter/pkg/collector"
	"go-logstash-exporter/pkg/server"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/spf13/cobra"
)

var (
	configFile string // 配置文件路径
)

// Run 是程序的主入口函数
func Run(cmd *cobra.Command, args []string) {
	// 必须指定配置文件
	if configFile == "" {
		fmt.Fprintf(os.Stderr, "错误: 必须通过 --config.file 指定配置文件\n")
		os.Exit(1)
	}

	var endpoints []string
	var bindAddress string = ":8080" // 默认监听地址，当配置文件中未指定时使用

	// 从配置文件读取
	config, err := server.LoadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置文件失败: %v\n", err)
		os.Exit(1)
	}
	endpoints = config.Endpoints
	if config.Web.ListenAddress != "" {
		bindAddress = config.Web.ListenAddress
	}

	// 注册系统信息收集器
	prometheus.MustRegister(collectors.NewBuildInfoCollector())

	// 为每个 endpoint 创建一个收集器
	for _, endpoint := range endpoints {
		endpoint = strings.TrimSpace(endpoint)
		if endpoint == "" {
			continue
		}

		// 创建并注册 Logstash 收集器
		logstashCollector, err := collector.New(endpoint)
		if err != nil {
			fmt.Fprintf(os.Stderr, "创建收集器失败 [%s]: %v\n", endpoint, err)
			continue
		}
		prometheus.MustRegister(logstashCollector)
		fmt.Printf("添加 Logstash 实例: %s\n", endpoint)
	}

	// 创建并启动 HTTP 服务器
	srv := server.New(bindAddress)
	srv.SetupRoutes()

	fmt.Printf("启动 Logstash 指标采集器，监听地址: %s\n", bindAddress)
	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "启动服务失败: %v\n", err)
		os.Exit(1)
	}
}
