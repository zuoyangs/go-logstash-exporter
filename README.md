# Logstash Exporter

Prometheus exporter for Logstash metrics.

## 功能特性

- 支持多个 Logstash 实例监控
- 支持配置文件方式配置
- 提供丰富的监控指标（JVM、进程、Pipeline 等）
- 支持 Logstash 5.x 和 6.x+ 版本

## 程序架构

### 目录结构

```
go-logstash-exporter/
├── cmd/
│   └── logstash_exporter.go  # 程序入口
├── pkg/
│   ├── collector/            # 指标收集器
│   │   ├── collector.go      # 主收集器
│   │   ├── nodestats_api.go  # 节点统计 API
│   │   ├── nodestats_collector.go  # 节点统计收集器
│   │   ├── nodeinfo_api.go   # 节点信息 API
│   │   └── nodeinfo_collector.go   # 节点信息收集器
│   └── server/               # HTTP 服务器
│       ├── server.go         # 服务器实现
│       └── config.go         # 配置处理
├── config.yaml               # 配置文件
└── main.go                   # 程序主入口
```

### 调用流程

1. **启动流程**:
```
main()
  → cmd.Execute()
    → cmd.Run() 
      → server.LoadConfig()        # 加载配置
      → collector.New()            # 创建收集器
        → NewNodeStatsCollector()  # 创建节点统计收集器
        → NewNodeInfoCollector()   # 创建节点信息收集器
      → server.New()              # 创建 HTTP 服务器
        → SetupRoutes()           # 设置路由
        → Start()                 # 启动服务
```

2. **指标收集流程**:
```
HTTP 请求 
  → server.Engine
    → prometheus.Handler
      → LogstashCollector.Collect
        → NodeStatsCollector.Collect
          → NodeStats()
            → HTTP GET /_node/stats
        → NodeInfoCollector.Collect
          → NodeInfo()
            → HTTP GET /_node/info
```

3. **数据流向**:
```
Logstash API → HTTP Client → Collector → Prometheus Metrics → HTTP Response
```

### 主要组件职责

1. **main.go 和 cmd/logstash_exporter.go**:
   - 程序入口点
   - 命令行参数处理
   - 配置文件加载
   - 服务启动管理

2. **pkg/collector**:
   - 实现指标收集逻辑
   - 管理多个收集器
   - 与 Logstash API 交互
   - 数据格式转换

3. **pkg/server**:
   - HTTP 服务器实现
   - 路由配置管理
   - 指标暴露接口

## 配置说明

### 命令行参数

```bash
logstash_exporter [flags]

Flags:
  --config.file string          配置文件路径（支持 YAML、JSON、TOML 等格式）
```

### 配置文件示例 (YAML)

```yaml
endpoints:
  - http://localhost:9600
  - http://logstash-02:9600
  - http://logstash-03:9600

web:
  listen_address: ":9198"
```

## 监控指标

### 核心指标类别

1. **JVM 相关指标**:
   - 线程数量统计
   - 内存使用情况
   - 垃圾回收统计

2. **进程指标**:
   - 文件描述符使用
   - CPU 使用情况
   - 虚拟内存使用

3. **Pipeline 指标**:
   - 事件处理统计
   - 插件性能指标
   - 队列状态监控

4. **插件性能指标**:
   - Input 插件指标
   - Filter 插件指标
   - Output 插件指标

## 使用示例

### 使用配置文件启动:

```bash
./logstash_exporter --config.file=config.yaml
```

### 访问指标接口:

```bash
curl http://localhost:9198/metrics
```

## 多实例监控

go-logstash-exporter 支持监控多个 Logstash 实例。在配置文件中列出所有要监控的 Logstash 实例端点:

```yaml
endpoints:
  - http://logstash1:9600
  - http://logstash2:9600
  - http://logstash3:9600
```

每个实例的指标都会带有 `instance` 标签，可以通过该标签区分不同实例的指标:

```promql
# 查看所有实例的堆内存使用情况
logstash_stats_mem_heap_used_bytes

# 查看特定实例的堆内存使用情况
logstash_stats_mem_heap_used_bytes{instance="logstash1:9600"}
```

## 开发说明

项目使用 Go 模块进行依赖管理。主要依赖：

- github.com/prometheus/client_golang：Prometheus 客户端库
- github.com/spf13/cobra：命令行工具库
- github.com/spf13/viper：配置管理库
- github.com/gin-gonic/gin：HTTP 框架

## License

[MIT License](LICENSE)