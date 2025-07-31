// collector 包提供了 Logstash 相关的数据收集功能
package collector

// NodeInfoResponse 定义了 Logstash 节点信息的响应结构，包含了节点的基本信息、管道配置、操作系统信息和 JVM 信息
type NodeInfoResponse struct {
	Host        string `json:"host"`         // Host 表示 Logstash 实例的主机名
	Version     string `json:"version"`      // Version 表示 Logstash 的版本号
	HTTPAddress string `json:"http_address"` // HTTPAddress 表示 Logstash HTTP API 的监听地址
	ID          string `json:"id"`           // ID 表示节点的唯一标识符
	Name        string `json:"name"`         // Name 表示节点的名称
	EphemeralID string `json:"ephemeral_id"` // EphemeralID 表示节点的临时标识符
	Status      string `json:"status"`       // Status 表示节点状态
	Snapshot    bool   `json:"snapshot"`     // Snapshot 表示是否为快照版本

	Pipeline struct {
		Workers               int  `json:"workers"`                 // Workers 表示处理事件的工作线程数
		BatchSize             int  `json:"batch_size"`              // BatchSize 表示每批处理的事件数量
		BatchDelay            int  `json:"batch_delay"`             // BatchDelay 表示批处理的延迟时间
		ConfigReloadAutomatic bool `json:"config_reload_automatic"` // ConfigReloadAutomatic 表示是否启用自动重载配置
		ConfigReloadInterval  int  `json:"config_reload_interval"`  // ConfigReloadInterval 表示配置重载的时间间隔
	} `json:"pipeline"` // Pipeline 包含了 Logstash 管道的配置信息

	Pipelines map[string]struct {
		EphemeralID           string `json:"ephemeral_id"`            // Pipeline 临时标识符
		Hash                  string `json:"hash"`                    // Pipeline 配置哈希值
		Workers               int    `json:"workers"`                 // Workers 表示处理事件的工作线程数
		BatchSize             int    `json:"batch_size"`              // BatchSize 表示每批处理的事件数量
		BatchDelay            int    `json:"batch_delay"`             // BatchDelay 表示批处理的延迟时间
		ConfigReloadAutomatic bool   `json:"config_reload_automatic"` // ConfigReloadAutomatic 表示是否启用自动重载配置
		ConfigReloadInterval  int64  `json:"config_reload_interval"`  // ConfigReloadInterval 表示配置重载的时间间隔（纳秒）
		DeadLetterQueueEnabled bool  `json:"dead_letter_queue_enabled"` // DeadLetterQueueEnabled 表示是否启用死信队列
	} `json:"pipelines"` // Pipelines 包含了 Logstash 多管道的配置信息

	Os struct {
		Name                string `json:"name"`                 // Name 表示操作系统名称
		Arch                string `json:"arch"`                 // Arch 表示系统架构
		Version             string `json:"version"`              // Version 表示操作系统版本
		AvailableProcessors int    `json:"available_processors"` // AvailableProcessors 表示可用的处理器数量
	} `json:"os"` // Os 包含了运行 Logstash 的操作系统信息

	Jvm struct {
		Pid               int    `json:"pid"`                  // Pid 表示 JVM 进程 ID
		Version           string `json:"version"`              // Version 表示 Java 版本
		VMName            string `json:"vm_name"`              // VMName 表示虚拟机名称
		VMVersion         string `json:"vm_version"`           // VMVersion 表示虚拟机版本
		VMVendor          string `json:"vm_vendor"`            // VMVendor 表示虚拟机供应商
		StartTimeInMillis int64  `json:"start_time_in_millis"` // StartTimeInMillis 表示 JVM 启动时间（毫秒）
		Mem               struct {
			HeapInitInBytes    int `json:"heap_init_in_bytes"`     // HeapInitInBytes 表示堆内存初始大小（字节）
			HeapMaxInBytes     int `json:"heap_max_in_bytes"`      // HeapMaxInBytes 表示堆内存最大大小（字节）
			NonHeapInitInBytes int `json:"non_heap_init_in_bytes"` // NonHeapInitInBytes 表示非堆内存初始大小（字节）
			NonHeapMaxInBytes  int `json:"non_heap_max_in_bytes"`  // NonHeapMaxInBytes 表示非堆内存最大大小（字节）
		} `json:"mem"` // Mem 包含了 JVM 内存使用情况
		GcCollectors []string `json:"gc_collectors"` // GcCollectors 表示垃圾收集器列表
	} `json:"jvm"` // Jvm 包含了 Java 虚拟机的相关信息
}

// NodeInfo 函数用于获取 Logstash 节点的详细信息
// endpoint 参数指定 Logstash API 的基础 URL
// 返回节点信息响应和可能的错误
func NodeInfo(endpoint string) (NodeInfoResponse, error) {
	var response NodeInfoResponse

	handler := &HTTPHandler{
		Endpoint: endpoint + "/_node",
	}

	err := getMetrics(handler, &response)

	return response, err
}