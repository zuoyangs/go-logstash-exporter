package main

import (
	"os"

	"go-logstash-exporter/cmd"

	"github.com/spf13/cobra"
)

var (
	configFile string // 配置文件路径
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "logstash_exporter",
		Short: "Logstash 指标采集器",
		Long:  "Logstash 指标采集器可以从一个或多个 Logstash 实例收集监控指标。",
		Run:   cmd.Run,
	}

	rootCmd.Flags().StringVar(&configFile, "config.file", "", "配置文件路径")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
