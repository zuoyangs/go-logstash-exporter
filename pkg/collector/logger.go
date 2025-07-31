package collector

import (
	"fmt"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// logger 是全局日志记录器实例
var logger log.Logger

// init 初始化日志记录器
// 设置日志输出到标准错误输出
// 设置日志级别过滤器，允许 Info 级别及以上的日志
func init() {
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = level.NewFilter(logger, level.AllowInfo())
}

// Error 记录错误级别的日志
// keyvals 参数为键值对形式的日志内容
func Error(keyvals ...interface{}) {
	level.Error(logger).Log(keyvals...)
}

// Info 记录信息级别的日志
// keyvals 参数为键值对形式的日志内容
func Info(keyvals ...interface{}) {
	level.Info(logger).Log(keyvals...)
}

// Debug 记录调试级别的日志
// keyvals 参数为键值对形式的日志内容
func Debug(keyvals ...interface{}) {
	level.Debug(logger).Log(keyvals...)
}

// Fatal 记录致命错误日志并退出程序
// keyvals 参数为键值对形式的日志内容
func Fatal(keyvals ...interface{}) {
	level.Error(logger).Log(keyvals...)
	os.Exit(1)
}

// Errorf 记录格式化的错误日志
// format 为格式化字符串
// args 为格式化参数
func Errorf(format string, args ...interface{}) {
	Error("msg", fmt.Sprintf(format, args...))
}

// Infof 记录格式化的信息日志
// format 为格式化字符串
// args 为格式化参数
func Infof(format string, args ...interface{}) {
	Info("msg", fmt.Sprintf(format, args...))
}

// Debugf 记录格式化的调试日志
// format 为格式化字符串
// args 为格式化参数
func Debugf(format string, args ...interface{}) {
	Debug("msg", fmt.Sprintf(format, args...))
}

// Fatalf 记录格式化的致命错误日志并退出程序
// format 为格式化字符串
// args 为格式化参数
func Fatalf(format string, args ...interface{}) {
	Fatal("msg", fmt.Sprintf(format, args...))
}