package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

// NodeInfoCollector 节点信息收集器
type NodeInfoCollector struct {
	endpoint string // Logstash API 端点
	instance string // 实例标识

	NodeInfos *prometheus.Desc // 节点信息指标
	OsInfos   *prometheus.Desc // 操作系统信息指标
	JvmInfos  *prometheus.Desc // JVM 信息指标
}

// NewNodeInfoCollector 创建新的节点信息收集器
func NewNodeInfoCollector(logstashEndpoint string, instance string) (Collector, error) {
	const subsystem = "info"

	return &NodeInfoCollector{
		endpoint: logstashEndpoint,
		instance: instance,

		NodeInfos: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "node"),
			"A metric with a constant '1' value labeled by Logstash version.",
			[]string{"version", "instance"},
			nil,
		),

		OsInfos: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "os"),
			"A metric with a constant '1' value labeled by name, arch, version and available_processors to the OS running Logstash.",
			[]string{"name", "arch", "version", "available_processors", "instance"},
			nil,
		),

		JvmInfos: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "jvm"),
			"A metric with a constant '1' value labeled by name, version and vendor of the JVM running Logstash.",
			[]string{"name", "version", "vendor", "instance"},
			nil,
		),
	}, nil
}

// Collect 收集节点信息指标
func (c *NodeInfoCollector) Collect(ch chan<- prometheus.Metric) error {
	if _, err := c.collect(ch); err != nil {
		Errorf("Failed collecting info metrics: %v", err)
		return err
	}
	return nil
}

// collect 实际执行节点信息收集工作
func (c *NodeInfoCollector) collect(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	stats, err := NodeInfo(c.endpoint)
	if err != nil {
		return nil, err
	}

	ch <- prometheus.MustNewConstMetric(
		c.NodeInfos,
		prometheus.CounterValue,
		float64(1),
		stats.Version,
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.OsInfos,
		prometheus.CounterValue,
		float64(1),
		stats.Os.Name,
		stats.Os.Arch,
		stats.Os.Version,
		strconv.Itoa(stats.Os.AvailableProcessors),
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.JvmInfos,
		prometheus.CounterValue,
		float64(1),
		stats.Jvm.VMName,
		stats.Jvm.VMVersion,
		stats.Jvm.VMVendor,
		c.instance,
	)

	return nil, nil
}