package main

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/mydev/mydev-api/enums"
	"github.com/mydev/mydev-api/exception"
	"github.com/mydev/mydev-service/internal/config"
	"github.com/mydev/mydev-service/internal/handler"
	"github.com/mydev/mydev-service/internal/middleware"
)

func main() {
	// 加载配置
	cfg := config.GetConfig()

	// 初始化日志
	logger, err := config.InitLogger(cfg)
	if err != nil {
		slog.Error("Failed to initialize logger", "error", err)
		os.Exit(1)
	}
	logger.Info("Logger initialized", "level", cfg.LogLevel)

	// 创建 Fiber 应用
	app := fiber.New(fiber.Config{
		AppName:               "MyDev Service",
		DisableStartupMessage: false,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Fiber 默认错误处理，返回统一格式
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(map[string]interface{}{
				"code": int(enums.UnknownError),
				"msg":  "Internal Server Error.",
				"data": []interface{}{},
			})
		},
	})

	// 注册中间件（注意顺序）
	app.Use(middleware.TraceIdGenerator())      // 1. 首先生成/获取 traceId
	app.Use(middleware.TraceIdExtractor())       // 2. 从 header 提取 traceId
	app.Use(middleware.TraceIdPropagator())      // 3. 传播 traceId 到响应
	app.Use(middleware.Recovery())               // 4. 异常恢复
	app.Use(fiberLogger.New(fiberLogger.Config{            // 5. 请求日志
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))

	// 创建处理器
	healthHandler := handler.NewHealthHandler()
	divideHandler := handler.NewDivideHandler()

	// 注册路由
	api := app.Group("/srv/v1")
	api.Post("/hc", healthHandler.Handle)
	api.Post("/divide", divideHandler.Handle)

	// 启动服务
	logger.Info("Server starting", "port", cfg.Port)
	if err := app.Listen(cfg.Port); err != nil {
		logger.Error("Server error", "error", err)
		panic(exception.NewServiceException(int(enums.UnknownError), err.Error()))
	}
}
