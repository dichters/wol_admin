package response

import (
	"io"
	"os"
	"sync"
)

// MultiWriter 多目标写入器，根据日志级别决定写入目标
type MultiWriter struct {
	stdout      io.Writer
	serviceFile io.Writer
	errorFile   io.Writer
}

// NewMultiWriter 创建多目标写入器
func NewMultiWriter(stdout, serviceFile, errorFile io.Writer) *MultiWriter {
	return &MultiWriter{
		stdout:      stdout,
		serviceFile: serviceFile,
		errorFile:   errorFile,
	}
}

// Write 实现 io.Writer 接口
func (m *MultiWriter) Write(p []byte) (n int, err error) {
	// 写入控制台
	m.stdout.Write(p)

	// 写入 service.log
	m.serviceFile.Write(p)

	// 简单判断是否为 ERROR 级别（JSON 格式中包含 "level":"ERROR"）
	if contains(p, `"level":"ERROR"`) || contains(p, `"level":"error"`) {
		m.errorFile.Write(p)
	}

	return len(p), nil
}

// contains 检查字节切片是否包含子串
func contains(data []byte, substr string) bool {
	for i := 0; i <= len(data)-len(substr); i++ {
		if string(data[i:i+len(substr)]) == substr {
			return true
		}
	}
	return false
}

// Sync 实现 syncer 接口（用于日志轮转）
func (m *MultiWriter) Sync() error {
	// TODO: 实现日志轮转逻辑
	return nil
}

// DailyRotator 每日日志轮转
type DailyRotator struct {
	mu          sync.Mutex
	currentDate string
	serviceFile *os.File
	errorFile   *os.File
	basePath    string
}

// NewDailyRotator 创建每日轮转器
func NewDailyRotator(basePath string) *DailyRotator {
	return &DailyRotator{
		basePath: basePath,
	}
}
