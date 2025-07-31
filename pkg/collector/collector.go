package collector

import (
	"net/url"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Namespace 定义了 Prometheus 指标的命名空间
const (
	Namespace = "logstash"
)

// 定义全局指标：抓取持续时间统计
var (
	scrapeDurations = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: Namespace,
			Subsystem: "exporter",
			Name:      "scrape_duration_seconds",
			Help:      "logstash_exporter: 抓取任务的持续时间统计。",
		},
		[]string{"collector", "result", "instance"},
	)
)

// Collector 接口定义了指标收集器的基本行为
type Collector interface {
	// Collect 方法用于收集指标并通过 channel 发送
	Collect(ch chan<- prometheus.Metric) error
}

// LogstashCollector 是主收集器，负责管理所有子收集器
type LogstashCollector struct {
	collectors map[string]Collector // 子收集器映射表
	endpoint   string               // Logstash API 端点
	instance   string               // Logstash 实例标识
}

// New 创建一个新的 LogstashCollector 实例
func New(endpoint string) (*LogstashCollector, error) {
	// 解析 endpoint URL 获取实例标识
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	instance := u.Host
	if instance == "" {
		instance = endpoint
	}

	// 创建节点统计信息收集器
	nodeStats, err := NewNodeStatsCollector(endpoint, instance)
	if err != nil {
		return nil, err
	}

	// 创建节点基本信息收集器
	nodeInfo, err := NewNodeInfoCollector(endpoint, instance)
	if err != nil {
		return nil, err
	}

	// 返回配置好的收集器实例
	return &LogstashCollector{
		endpoint: endpoint,
		instance: instance,
		collectors: map[string]Collector{
			"node": nodeStats, // 节点统计信息收集器
			"info": nodeInfo,  // 节点基本信息收集器
		},
	}, nil
}

// Describe 实现了 prometheus.Collector 接口，用于描述所有可能的指标
func (c *LogstashCollector) Describe(ch chan<- *prometheus.Desc) {
	scrapeDurations.Describe(ch)
}

// Collect 实现了 prometheus.Collector 接口，用于收集当前的指标值
func (c *LogstashCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(c.collectors))

	// 并发收集所有子收集器的指标
	for name, collector := range c.collectors {
		go func(name string, c Collector, instance string) {
			begin := time.Now()
			err := c.Collect(ch)
			duration := time.Since(begin)

			// 记录收集结果
			result := "success"
			if err != nil {
				result = "error"
			}

			// 更新抓取持续时间指标
			scrapeDurations.WithLabelValues(name, result, instance).Observe(duration.Seconds())
			wg.Done()
		}(name, collector, c.instance)
	}

	// 等待所有收集器完成
	wg.Wait()

	// 收集抓取持续时间指标
	scrapeDurations.Collect(ch)
}