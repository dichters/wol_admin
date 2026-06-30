package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// TraceIdPropagator MDC ID 传播中间件
// 将 traceId 和 spanId 传播到响应头
func TraceIdPropagator() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()

		// 将 traceId 和 spanId 添加到响应头
		if traceId := c.Locals(TraceIdKey); traceId != nil {
			c.Set(HeaderTrace, ToString(traceId))
		}
		if spanId := c.Locals(SpanIdKey); spanId != nil {
			c.Set(HeaderSpan, ToString(spanId))
		}

		return err
	}
}
