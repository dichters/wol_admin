package config

import (
	"log/slog"
	"os"
	"sync"

	"github.com/mydev/mydev-api/response"
)

// Config 应用配置
type Config struct {
	Port     string
	LogLevel slog.Level
}

var (
	instance *Config
	once     sync.Once
)

// GetConfig 获取配置单例
func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{
			Port:     ":8030",
			LogLevel: slog.LevelInfo,
		}
	})
	return instance
}

// InitLogger 初始化日志
func InitLogger(cfg *Config) (*slog.Logger, error) {
	// 确保日志目录存在
	if err := os.MkdirAll("log", 0755); err != nil {
		return nil, err
	}

	// 创建 service.log 文件
	serviceFile, err := os.OpenFile("log/service.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	// 创建 error.log 文件
	errorFile, err := os.OpenFile("log/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		serviceFile.Close()
		return nil, err
	}

	// 创建多写入器
	multiWriter := response.NewMultiWriter(
		os.Stdout,   // 控制台输出
		serviceFile, // service.log 输出所有日志
		errorFile,   // error.log 输出 ERROR 及以上日志
	)

	logger := slog.New(slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))
	slog.SetDefault(logger)

	return logger, nil
}
