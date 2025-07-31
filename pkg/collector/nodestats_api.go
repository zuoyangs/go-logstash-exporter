package collector

// Pipeline 结构体定义了 Logstash pipeline 的所有监控指标
type Pipeline struct {
	// Events 记录整个 pipeline 的事件处理统计
	Events struct {
		DurationInMillis int `json:"duration_in_millis"` // pipeline 处理事件的总耗时（毫秒）
		In               int `json:"in"`                 // 进入 pipeline 的事件总数
		Filtered         int `json:"filtered"`           // 经过过滤器处理的事件数
		Out              int `json:"out"`                // 从 pipeline 输出的事件总数
	} `json:"events"`

	// Plugins 包含所有插件（inputs、filters、outputs）的性能指标
	Plugins struct {
		// Inputs 记录所有输入插件的性能指标
		Inputs []struct {
			ID     string `json:"id"` // 插件实例的唯一标识符
			Events struct {
				In  int `json:"in"`  // 插件接收到的原始事件数
				Out int `json:"out"` // 插件成功处理并输出的事件数
			} `json:"events"`
			Name string `json:"name"` // 插件的名称（如 beats、file、kafka 等）
		} `json:"inputs,omitempty"`

		// Filters 记录所有过滤器插件的性能指标
		Filters []struct {
			ID     string `json:"id"` // 过滤器实例的唯一标识符
			Events struct {
				DurationInMillis int `json:"duration_in_millis"` // 过滤器处理事件的总耗时
				In               int `json:"in"`                 // 进入过滤器的事件数
				Out              int `json:"out"`                // 过滤器处理后输出的事件数
			} `json:"events,omitempty"`
			Name             string `json:"name"`     // 过滤器名称（如 grok、mutate、date 等）
			Matches          int    `json:"matches"`  // 匹配成功的事件数
			Failures         int    `json:"failures"` // 处理失败的事件数
			PatternsPerField struct {
				CapturedRequestHeaders int `json:"captured_request_headers"` // 捕获的请求头数量
			} `json:"patterns_per_field,omitempty"`
			Formats int `json:"formats,omitempty"` // 支持的格式数量
		} `json:"filters"`

		// Outputs 记录所有输出插件的性能指标
		Outputs []struct {
			ID     string `json:"id"` // 输出插件实例的唯一标识符
			Events struct {
				In  int `json:"in"`  // 进入输出插件的事件数
				Out int `json:"out"` // 成功发送的事件数
			} `json:"events"`
			Name string `json:"name"` // 输出插件名称（如 elasticsearch、kafka、file 等）
		} `json:"outputs"`
	} `json:"plugins"`

	// Reloads 记录 pipeline 配置重载的统计信息
	Reloads struct {
		LastError            interface{} `json:"last_error"`             // 最后一次重载错误信息
		Successes            int         `json:"successes"`              // 成功重载次数
		LastSuccessTimestamp interface{} `json:"last_success_timestamp"` // 最后一次成功重载时间
		LastFailureTimestamp interface{} `json:"last_failure_timestamp"` // 最后一次失败重载时间
		Failures             int         `json:"failures"`               // 重载失败次数
	} `json:"reloads"`

	// Queue 记录队列相关的性能指标
	Queue struct {
		Events   int    `json:"events"` // 当前队列中的事件数量
		Type     string `json:"type"`   // 队列类型（memory 或 persisted）
		Capacity struct {
			PageCapacityInBytes int   `json:"page_capacity_in_bytes"`  // 每个队列页的容量（字节）
			MaxQueueSizeInBytes int64 `json:"max_queue_size_in_bytes"` // 队列最大容量（字节）
			MaxUnreadEvents     int   `json:"max_unread_events"`       // 最大未读事件数
		} `json:"capacity"`
		Data struct {
			Path             string `json:"path"`                // 持久化队列的存储路径
			FreeSpaceInBytes int64  `json:"free_space_in_bytes"` // 剩余可用空间（字节）
			StorageType      string `json:"storage_type"`        // 存储类型
		} `json:"data"`
	} `json:"queue"`

	// DeadLetterQueue 记录死信队列的统计信息
	DeadLetterQueue struct {
		QueueSizeInBytes int `json:"queue_size_in_bytes"` // 死信队列大小（字节）
	} `json:"dead_letter_queue"`
}

// NodeStatsResponse 定义了从 Logstash 节点获取的所有统计信息
type NodeStatsResponse struct {
	Host        string `json:"host"`         // Logstash 节点主机名
	Version     string `json:"version"`      // Logstash 版本号
	HTTPAddress string `json:"http_address"` // HTTP API 监听地址

	// JVM 运行时统计信息
	Jvm struct {
		// 线程统计
		Threads struct {
			Count     int `json:"count"`      // 当前活动线程数
			PeakCount int `json:"peak_count"` // 峰值线程数
		} `json:"threads"`

		// 内存使用统计
		Mem struct {
			HeapUsedInBytes         int `json:"heap_used_in_bytes"`          // 已使用的堆内存（字节）
			HeapUsedPercent         int `json:"heap_used_percent"`           // 堆内存使用百分比
			HeapCommittedInBytes    int `json:"heap_committed_in_bytes"`     // 已提交的堆内存（字节）
			HeapMaxInBytes          int `json:"heap_max_in_bytes"`           // 最大堆内存限制（字节）
			NonHeapUsedInBytes      int `json:"non_heap_used_in_bytes"`      // 已使用的非堆内存（字节）
			NonHeapCommittedInBytes int `json:"non_heap_committed_in_bytes"` // 已提交的非堆内存（字节）

			// 内存池统计（包括新生代、老年代和幸存区）
			Pools struct {
				// 幸存区内存池
				Survivor struct {
					PeakUsedInBytes  int `json:"peak_used_in_bytes"` // 峰值使用内存（字节）
					UsedInBytes      int `json:"used_in_bytes"`      // 当前使用内存（字节）
					PeakMaxInBytes   int `json:"peak_max_in_bytes"`  // 历史最大内存限制（字节）
					MaxInBytes       int `json:"max_in_bytes"`       // 当前最大内存限制（字节）
					CommittedInBytes int `json:"committed_in_bytes"` // 已提交内存（字节）
				} `json:"survivor"`

				// 老年代内存池
				Old struct {
					PeakUsedInBytes  int `json:"peak_used_in_bytes"` // 峰值使用内存（字节）
					UsedInBytes      int `json:"used_in_bytes"`      // 当前使用内存（字节）
					PeakMaxInBytes   int `json:"peak_max_in_bytes"`  // 历史最大内存限制（字节）
					MaxInBytes       int `json:"max_in_bytes"`       // 当前最大内存限制（字节）
					CommittedInBytes int `json:"committed_in_bytes"` // 已提交内存（字节）
				} `json:"old"`

				// 新生代内存池
				Young struct {
					PeakUsedInBytes  int `json:"peak_used_in_bytes"` // 峰值使用内存（字节）
					UsedInBytes      int `json:"used_in_bytes"`      // 当前使用内存（字节）
					PeakMaxInBytes   int `json:"peak_max_in_bytes"`  // 历史最大内存限制（字节）
					MaxInBytes       int `json:"max_in_bytes"`       // 当前最大内存限制（字节）
					CommittedInBytes int `json:"committed_in_bytes"` // 已提交内存（字节）
				} `json:"young"`
			} `json:"pools"`
		} `json:"mem"`

		// 垃圾回收统计
		Gc struct {
			Collectors struct {
				// 老年代垃圾回收器
				Old struct {
					CollectionTimeInMillis int `json:"collection_time_in_millis"` // 垃圾回收总时间（毫秒）
					CollectionCount        int `json:"collection_count"`          // 垃圾回收次数
				} `json:"old"`

				// 新生代垃圾回收器
				Young struct {
					CollectionTimeInMillis int `json:"collection_time_in_millis"` // 垃圾回收总时间（毫秒）
					CollectionCount        int `json:"collection_count"`          // 垃圾回收次数
				} `json:"young"`
			} `json:"collectors"`
		} `json:"gc"`
	} `json:"jvm"`

	// 进程统计信息
	Process struct {
		OpenFileDescriptors     int `json:"open_file_descriptors"`      // 当前打开的文件描述符数量
		PeakOpenFileDescriptors int `json:"peak_open_file_descriptors"` // 峰值文件描述符数量
		MaxFileDescriptors      int `json:"max_file_descriptors"`       // 最大文件描述符限制

		// 进程内存统计
		Mem struct {
			TotalVirtualInBytes int64 `json:"total_virtual_in_bytes"` // 总虚拟内存使用量（字节）
		} `json:"mem"`

		// CPU 使用统计
		CPU struct {
			TotalInMillis int64 `json:"total_in_millis"` // CPU 使用总时间（毫秒）
			Percent       int   `json:"percent"`         // CPU 使用百分比
		} `json:"cpu"`
	} `json:"process"`

	Pipeline  Pipeline            `json:"pipeline"`  // Logstash 5.x 的单 pipeline 统计
	Pipelines map[string]Pipeline `json:"pipelines"` // Logstash 6.x+ 的多 pipeline 统计
}

// NodeStats 函数从 Logstash 节点的 /_node/stats API 获取统计信息
func NodeStats(endpoint string) (NodeStatsResponse, error) {
	var response NodeStatsResponse

	// 创建 HTTP 处理器
	handler := &HTTPHandler{
		Endpoint: endpoint + "/_node/stats",
	}

	// 获取并解析统计数据
	err := getMetrics(handler, &response)

	return response, err
}
