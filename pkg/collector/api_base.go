package collector

import (
	"encoding/json"
	"net/http"
)

// HTTPHandler HTTP处理器结构体
type HTTPHandler struct {
	Endpoint string // 端点URL
}

// Get 发送HTTP GET请求并返回响应
func (h *HTTPHandler) Get() (http.Response, error) {
	response, err := http.Get(h.Endpoint)
	if err != nil {
		return http.Response{}, err
	}

	return *response, nil
}

// HTTPHandlerInterface HTTP处理器接口
type HTTPHandlerInterface interface {
	Get() (http.Response, error)
}

// getMetrics 从HTTP处理器获取指标数据并解析到目标结构体
func getMetrics(h HTTPHandlerInterface, target interface{}) error {
	response, err := h.Get()
	if err != nil {
		Errorf("无法获取指标: %s", err)
		return nil
	}

	defer func() {
		err = response.Body.Close()
		if err != nil {
			Errorf("无法关闭响应体: %v", err)
		}
	}()

	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		Errorf("无法解析Logstash响应json: %s", err)
	}

	return nil
}