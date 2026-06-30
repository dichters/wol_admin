package client

import (
	"fmt"
)

// Example SDK 使用示例
func Example() {
	// 创建客户端
	client := NewClient("http://localhost:8080")

	// 健康检查
	result, err := client.HealthCheck("trace-123", "span-456")
	if err != nil {
		fmt.Printf("HealthCheck error: %v\n", err)
		return
	}
	fmt.Printf("HealthCheck result: code=%d, msg=%s\n", result.Code, result.Msg)

	// 除法运算
	result, err = client.Divide(10, 2, "trace-789", "span-012")
	if err != nil {
		fmt.Printf("Divide error: %v\n", err)
		return
	}
	fmt.Printf("Divide result: code=%d, msg=%s, data=%v\n", result.Code, result.Msg, result.Data)
}
