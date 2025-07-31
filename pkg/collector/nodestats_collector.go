package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

// NodeStatsCollector 负责收集 Logstash 节点的统计信息
type NodeStatsCollector struct {
	endpoint string // Logstash API 端点
	instance string // Logstash 实例标识

	// JVM 相关指标
	JvmThreadsCount     *prometheus.Desc // JVM 线程数
	JvmThreadsPeakCount *prometheus.Desc // JVM 峰值线程数

	// 内存相关指标
	MemHeapUsedInBytes         *prometheus.Desc // 堆内存使用量
	MemHeapCommittedInBytes    *prometheus.Desc // 堆内存提交量
	MemHeapMaxInBytes          *prometheus.Desc // 堆内存最大值
	MemNonHeapUsedInBytes      *prometheus.Desc // 非堆内存使用量
	MemNonHeapCommittedInBytes *prometheus.Desc // 非堆内存提交量

	// 内存池相关指标
	MemPoolPeakUsedInBytes  *prometheus.Desc // 内存池峰值使用量
	MemPoolUsedInBytes      *prometheus.Desc // 内存池当前使用量
	MemPoolPeakMaxInBytes   *prometheus.Desc // 内存池峰值最大值
	MemPoolMaxInBytes       *prometheus.Desc // 内存池最大值
	MemPoolCommittedInBytes *prometheus.Desc // 内存池提交量

	// GC 相关指标
	GCCollectionTimeInMillis *prometheus.Desc // GC 收集时间
	GCCollectionCount        *prometheus.Desc // GC 收集次数

	// 进程相关指标
	ProcessOpenFileDescriptors    *prometheus.Desc // 打开的文件描述符数量
	ProcessMaxFileDescriptors     *prometheus.Desc // 最大文件描述符限制
	ProcessMemTotalVirtualInBytes *prometheus.Desc // 虚拟内存总量
	ProcessCPUTotalInMillis       *prometheus.Desc // CPU 使用时间

	// Pipeline 整体指标
	PipelineDuration       *prometheus.Desc // Pipeline 处理时间
	PipelineEventsIn       *prometheus.Desc // Pipeline 输入事件数
	PipelineEventsFiltered *prometheus.Desc // Pipeline 过滤事件数
	PipelineEventsOut      *prometheus.Desc // Pipeline 输出事件数

	// Pipeline 插件指标
	PipelinePluginEventsDuration *prometheus.Desc // 插件处理时间
	PipelinePluginEventsIn       *prometheus.Desc // 插件输入事件数
	PipelinePluginEventsOut      *prometheus.Desc // 插件输出事件数
	PipelinePluginMatches        *prometheus.Desc // 插件匹配次数
	PipelinePluginFailures       *prometheus.Desc // 插件失败次数

	// Pipeline 队列指标
	PipelineQueueEvents          *prometheus.Desc // 队列中的事件数
	PipelineQueuePageCapacity    *prometheus.Desc // 队列页容量
	PipelineQueueMaxQueueSize    *prometheus.Desc // 队列最大大小
	PipelineQueueMaxUnreadEvents *prometheus.Desc // 队列最大未读事件数

	// 死信队列指标
	PipelineDeadLetterQueueSizeInBytes *prometheus.Desc // 死信队列大小
}

// NewNodeStatsCollector 创建新的节点统计信息收集器
func NewNodeStatsCollector(logstashEndpoint string, instance string) (Collector, error) {
	const subsystem = "node"

	return &NodeStatsCollector{
		endpoint: logstashEndpoint,
		instance: instance,

		JvmThreadsCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "jvm_threads_count"),
			"jvm_threads_count",
			[]string{"instance"},
			nil,
		),

		JvmThreadsPeakCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "jvm_threads_peak_count"),
			"jvm_threads_peak_count",
			nil,
			nil,
		),

		MemHeapUsedInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "mem_heap_used_bytes"),
			"mem_heap_used_bytes",
			nil,
			nil,
		),

		MemHeapCommittedInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "mem_heap_committed_bytes"),
			"mem_heap_committed_bytes",
			nil,
			nil,
		),

		MemHeapMaxInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "mem_heap_max_bytes"),
			"mem_heap_max_bytes",
			nil,
			nil,
		),

		MemNonHeapUsedInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "mem_nonheap_used_bytes"),
			"mem_nonheap_used_bytes",
			nil,
			nil,
		),

		MemNonHeapCommittedInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "mem_nonheap_committed_bytes"),
			"mem_nonheap_committed_bytes",
			nil,
			nil,
		),

		MemPoolUsedInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "mem_pool_used_bytes"),
			"mem_pool_used_bytes",
			[]string{"pool"},
			nil,
		),

		MemPoolPeakUsedInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "mem_pool_peak_used_bytes"),
			"mem_pool_peak_used_bytes",
			[]string{"pool"},
			nil,
		),

		MemPoolPeakMaxInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "mem_pool_peak_max_bytes"),
			"mem_pool_peak_max_bytes",
			[]string{"pool"},
			nil,
		),

		MemPoolMaxInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "mem_pool_max_bytes"),
			"mem_pool_max_bytes",
			[]string{"pool"},
			nil,
		),

		MemPoolCommittedInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "mem_pool_committed_bytes"),
			"mem_pool_committed_bytes",
			[]string{"pool"},
			nil,
		),

		GCCollectionTimeInMillis: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "gc_collection_duration_seconds_total"),
			"gc_collection_duration_seconds_total",
			[]string{"collector"},
			nil,
		),

		GCCollectionCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "gc_collection_total"),
			"gc_collection_total",
			[]string{"collector"},
			nil,
		),

		ProcessOpenFileDescriptors: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "process_open_filedescriptors"),
			"process_open_filedescriptors",
			nil,
			nil,
		),

		ProcessMaxFileDescriptors: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "process_max_filedescriptors"),
			"process_max_filedescriptors",
			nil,
			nil,
		),

		ProcessMemTotalVirtualInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "process_mem_total_virtual_bytes"),
			"process_mem_total_virtual_bytes",
			nil,
			nil,
		),

		ProcessCPUTotalInMillis: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "process_cpu_total_seconds_total"),
			"process_cpu_total_seconds_total",
			nil,
			nil,
		),

		PipelineDuration: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "pipeline_duration_seconds_total"),
			"pipeline_duration_seconds_total",
			[]string{"pipeline"},
			nil,
		),

		PipelineEventsIn: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "pipeline_events_in_total"),
			"pipeline_events_in_total",
			[]string{"pipeline"},
			nil,
		),

		PipelineEventsFiltered: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "pipeline_events_filtered_total"),
			"pipeline_events_filtered_total",
			[]string{"pipeline"},
			nil,
		),

		PipelineEventsOut: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "pipeline_events_out_total"),
			"pipeline_events_out_total",
			[]string{"pipeline"},
			nil,
		),

		PipelinePluginEventsDuration: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "plugin_duration_seconds_total"),
			"plugin_duration_seconds",
			[]string{"pipeline", "plugin", "plugin_id", "plugin_type"},
			nil,
		),

		PipelinePluginEventsIn: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "plugin_events_in_total"),
			"plugin_events_in",
			[]string{"pipeline", "plugin", "plugin_id", "plugin_type"},
			nil,
		),

		PipelinePluginEventsOut: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "plugin_events_out_total"),
			"plugin_events_out",
			[]string{"pipeline", "plugin", "plugin_id", "plugin_type"},
			nil,
		),

		PipelinePluginMatches: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "plugin_matches_total"),
			"plugin_matches",
			[]string{"pipeline", "plugin", "plugin_id", "plugin_type"},
			nil,
		),

		PipelinePluginFailures: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "plugin_failures_total"),
			"plugin_failures",
			[]string{"pipeline", "plugin", "plugin_id", "plugin_type"},
			nil,
		),

		PipelineQueueEvents: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "queue_events"),
			"queue_events",
			[]string{"pipeline"},
			nil,
		),

		PipelineQueuePageCapacity: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "queue_page_capacity_bytes"),
			"queue_page_capacity_bytes",
			[]string{"pipeline"},
			nil,
		),

		PipelineQueueMaxQueueSize: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "queue_max_size_bytes"),
			"queue_max_size_bytes",
			[]string{"pipeline"},
			nil,
		),

		PipelineQueueMaxUnreadEvents: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "queue_max_unread_events"),
			"queue_max_unread_events",
			[]string{"pipeline"},
			nil,
		),

		PipelineDeadLetterQueueSizeInBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "dead_letter_queue_size_bytes"),
			"dead_letter_queue_size_bytes",
			[]string{"pipeline"},
			nil,
		),
	}, nil
}

// Collect 入口方法，负责错误处理和调用分发；
func (c *NodeStatsCollector) Collect(ch chan<- prometheus.Metric) error {
	if _, err := c.collect(ch); err != nil {
		Errorf("Failed collecting node stats metrics: %v", err)
		return err
	}
	return nil
}

// collect 方法实现了实际的指标收集逻辑
func (c *NodeStatsCollector) collect(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	stats, err := NodeStats(c.endpoint)
	if err != nil {
		return nil, err
	}

	ch <- prometheus.MustNewConstMetric(
		c.JvmThreadsCount,
		prometheus.GaugeValue,
		float64(stats.Jvm.Threads.Count),
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.JvmThreadsPeakCount,
		prometheus.GaugeValue,
		float64(stats.Jvm.Threads.PeakCount),
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemHeapUsedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.HeapUsedInBytes),
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemHeapMaxInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.HeapMaxInBytes),
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemHeapCommittedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.HeapCommittedInBytes),
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemNonHeapUsedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.NonHeapUsedInBytes),
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemNonHeapCommittedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.NonHeapCommittedInBytes),
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolPeakUsedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Old.PeakUsedInBytes),
		"old",
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolUsedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Old.UsedInBytes),
		"old",
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolPeakMaxInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Old.PeakMaxInBytes),
		"old",
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolMaxInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Old.MaxInBytes),
		"old",
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolCommittedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Old.CommittedInBytes),
		"old",
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolPeakUsedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Old.PeakUsedInBytes),
		"young",
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolUsedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Young.UsedInBytes),
		"young",
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolPeakMaxInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Old.PeakMaxInBytes),
		"young",
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolMaxInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Young.MaxInBytes),
		"young",
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolCommittedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Young.CommittedInBytes),
		"young",
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolPeakUsedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Old.PeakUsedInBytes),
		"survivor",
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolUsedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Survivor.UsedInBytes),
		"survivor",
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolPeakMaxInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Old.PeakMaxInBytes),
		"survivor",
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolMaxInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Survivor.MaxInBytes),
		"survivor",
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.MemPoolCommittedInBytes,
		prometheus.GaugeValue,
		float64(stats.Jvm.Mem.Pools.Survivor.CommittedInBytes),
		"survivor",
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.GCCollectionTimeInMillis,
		prometheus.CounterValue,
		float64(stats.Jvm.Gc.Collectors.Old.CollectionTimeInMillis),
		"old",
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.GCCollectionCount,
		prometheus.GaugeValue,
		float64(stats.Jvm.Gc.Collectors.Old.CollectionCount),
		"old",
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.GCCollectionTimeInMillis,
		prometheus.CounterValue,
		float64(stats.Jvm.Gc.Collectors.Young.CollectionTimeInMillis),
		"young",
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.GCCollectionCount,
		prometheus.GaugeValue,
		float64(stats.Jvm.Gc.Collectors.Young.CollectionCount),
		"young",
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.ProcessOpenFileDescriptors,
		prometheus.GaugeValue,
		float64(stats.Process.OpenFileDescriptors),
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.ProcessMaxFileDescriptors,
		prometheus.GaugeValue,
		float64(stats.Process.MaxFileDescriptors),
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.ProcessMemTotalVirtualInBytes,
		prometheus.GaugeValue,
		float64(stats.Process.Mem.TotalVirtualInBytes),
		c.instance,
	)

	ch <- prometheus.MustNewConstMetric(
		c.ProcessCPUTotalInMillis,
		prometheus.CounterValue,
		float64(stats.Process.CPU.TotalInMillis/1000),
		c.instance,
	)

	// 直接使用 Pipelines，不再支持 Logstash 5.x
	pipelines := stats.Pipelines

	for pipelineID, pipeline := range pipelines {
		ch <- prometheus.MustNewConstMetric(
			c.PipelineDuration,
			prometheus.CounterValue,
			float64(pipeline.Events.DurationInMillis/1000),
			pipelineID,
			c.instance,
		)

		ch <- prometheus.MustNewConstMetric(
			c.PipelineEventsIn,
			prometheus.CounterValue,
			float64(pipeline.Events.In),
			pipelineID,
			c.instance,
		)

		ch <- prometheus.MustNewConstMetric(
			c.PipelineEventsFiltered,
			prometheus.CounterValue,
			float64(pipeline.Events.Filtered),
			pipelineID,
			c.instance,
		)

		ch <- prometheus.MustNewConstMetric(
			c.PipelineEventsOut,
			prometheus.CounterValue,
			float64(pipeline.Events.Out),
			pipelineID,
			c.instance,
		)

		// 收集 Input 插件指标
		for _, plugin := range pipeline.Plugins.Inputs {
			// 输入事件计数 (适配 Logstash 7.5.0: 使用 Out 而不是 In)
			ch <- prometheus.MustNewConstMetric(
				c.PipelinePluginEventsIn,
				prometheus.CounterValue,
				float64(plugin.Events.Out),
				pipelineID,
				plugin.Name,
				plugin.ID,
				"input",
				c.instance,
			)
			// 输出事件计数
			ch <- prometheus.MustNewConstMetric(
				c.PipelinePluginEventsOut,
				prometheus.CounterValue,
				float64(plugin.Events.Out),
				pipelineID,
				plugin.Name,
				plugin.ID,
				"input",
				c.instance,
			)
		}

		// 收集 Filter 插件指标
		for _, plugin := range pipeline.Plugins.Filters {
			// 处理时间
			ch <- prometheus.MustNewConstMetric(
				c.PipelinePluginEventsDuration,
				prometheus.CounterValue,
				float64(plugin.Events.DurationInMillis/1000),
				pipelineID,
				plugin.Name,
				plugin.ID,
				"filter",
				c.instance,
			)
			// 输入事件计数
			ch <- prometheus.MustNewConstMetric(
				c.PipelinePluginEventsIn,
				prometheus.CounterValue,
				float64(plugin.Events.In),
				pipelineID,
				plugin.Name,
				plugin.ID,
				"filter",
				c.instance,
			)
			// 输出事件计数
			ch <- prometheus.MustNewConstMetric(
				c.PipelinePluginEventsOut,
				prometheus.CounterValue,
				float64(plugin.Events.Out),
				pipelineID,
				plugin.Name,
				plugin.ID,
				"filter",
				c.instance,
			)
			// 匹配成功计数
			ch <- prometheus.MustNewConstMetric(
				c.PipelinePluginMatches,
				prometheus.CounterValue,
				float64(plugin.Matches),
				pipelineID,
				plugin.Name,
				plugin.ID,
				"filter",
				c.instance,
			)
			// 处理失败计数
			ch <- prometheus.MustNewConstMetric(
				c.PipelinePluginFailures,
				prometheus.CounterValue,
				float64(plugin.Failures),
				pipelineID,
				plugin.Name,
				plugin.ID,
				"filter",
				c.instance,
			)
		}

		// 收集 Output 插件指标
		for _, plugin := range pipeline.Plugins.Outputs {
			// 输入事件计数
			ch <- prometheus.MustNewConstMetric(
				c.PipelinePluginEventsIn,
				prometheus.CounterValue,
				float64(plugin.Events.In),
				pipelineID,
				plugin.Name,
				plugin.ID,
				"output",
				c.instance,
			)
			// 输出事件计数
			ch <- prometheus.MustNewConstMetric(
				c.PipelinePluginEventsOut,
				prometheus.CounterValue,
				float64(plugin.Events.Out),
				pipelineID,
				plugin.Name,
				plugin.ID,
				"output",
				c.instance,
			)
		}

		if pipeline.Queue.Type != "memory" {
			ch <- prometheus.MustNewConstMetric(
				c.PipelineQueueEvents,
				prometheus.CounterValue,
				float64(pipeline.Queue.EventsCount),
				pipelineID,
				c.instance,
			)

			ch <- prometheus.MustNewConstMetric(
				c.PipelineQueuePageCapacity,
				prometheus.CounterValue,
				float64(pipeline.Queue.QueueSizeInBytes),
				pipelineID,
				c.instance,
			)

			ch <- prometheus.MustNewConstMetric(
				c.PipelineQueueMaxQueueSize,
				prometheus.CounterValue,
				float64(pipeline.Queue.MaxQueueSizeInBytes),
				pipelineID,
				c.instance,
			)
		}

	}

	return nil, nil
}
