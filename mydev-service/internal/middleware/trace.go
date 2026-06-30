package middleware

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TraceIdExtractor MDC ID 获取中间件
// 从请求头中提取 traceId 和 spanId
func TraceIdExtractor() fiber.Handler {
	return func(c *fiber.Ctx) error {
		traceId := c.Get(HeaderTrace)
		spanId := c.Get(HeaderSpan)

		if traceId != "" {
			c.Locals(TraceIdKey, traceId)
		}
		if spanId != "" {
			c.Locals(SpanIdKey, spanId)
		}

		return c.Next()
	}
}

// TraceIdGenerator MDC ID 生成中间件
// 如果请求头中没有 traceId 和 spanId，则生成新的
// 这个中间件应该最先执行
func TraceIdGenerator() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 检查是否已有 traceId
		if c.Get(HeaderTrace) == "" {
			traceId := uuid.New().String()
			c.Locals(TraceIdKey, traceId)
			c.Request().Header.Set(HeaderTrace, traceId)
		} else {
			c.Locals(TraceIdKey, c.Get(HeaderTrace))
		}

		// 检查是否已有 spanId
		if c.Get(HeaderSpan) == "" {
			spanId := uuid.New().String()
			c.Locals(SpanIdKey, spanId)
			c.Request().Header.Set(HeaderSpan, spanId)
		} else {
			c.Locals(SpanIdKey, c.Get(HeaderSpan))
		}

		// 添加到日志上下文
		logHandler := slog.Default()
		logHandler = logHandler.With(
			slog.String("traceId", c.Locals(TraceIdKey).(string)),
			slog.String("spanId", c.Locals(SpanIdKey).(string)),
		)
		slog.SetDefault(logHandler)

		return c.Next()
	}
}
