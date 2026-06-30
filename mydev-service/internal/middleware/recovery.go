package middleware

import (
	"log/slog"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/mydev/mydev-api/enums"
	"github.com/mydev/mydev-api/exception"
	"github.com/mydev/mydev-api/response"
)


// Recovery 异常恢复中间件
func Recovery() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// 获取日志上下文
				logHandler := recoveryLogContext(c)

				// 打印 panic 堆栈
				logHandler.Error("Panic recovered",
					slog.String("panic", ToString(r)),
					slog.String("stack", string(debug.Stack())),
				)

				// 返回未知错误响应
				c.Status(fiber.StatusInternalServerError).JSON(response.Error(
					int(enums.UnknownError),
					"Internal Server Error.",
				))
			}
		}()

		// 执行请求
		err := c.Next()

		// 处理错误
		if err != nil {
			return handleError(c, err)
		}

		return nil
	}
}

// handleError 处理错误
func handleError(c *fiber.Ctx, err error) error {
	logHandler := recoveryLogContext(c)

	// 将 traceId 和 spanId 添加到响应头
 propagateTraceIds(c)

	switch e := err.(type) {
	case *exception.BizException:
		// 业务错误：返回自定义 code 和 msg
		logHandler.Warn("Business error",
			slog.Int("code", e.Code),
			slog.String("msg", e.Msg),
		)
		return c.Status(fiber.StatusOK).JSON(response.Error(e.Code, e.Msg))

	case *exception.ServiceException:
		// 服务错误：返回自定义 code，但 msg 固定为 "Internal Server Error."
		logHandler.Error("Service error",
			slog.Int("code", e.Code),
			slog.String("original_msg", e.Msg),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(e.Code, "Internal Server Error."))

	case *exception.BaseException:
		// 基础异常：按未知错误处理
		logHandler.Error("BaseException error",
			slog.Int("code", e.Code),
			slog.String("msg", e.Msg),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(
			int(enums.UnknownError),
			"Internal Server Error.",
		))

	default:
		// 未知错误
		logHandler.Error("Unknown error",
			slog.String("error", err.Error()),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(
			int(enums.UnknownError),
			"Internal Server Error.",
		))
	}
}

// recoveryLogContext 获取带有 traceId 和 spanId 的日志上下文
func recoveryLogContext(c *fiber.Ctx) *slog.Logger {
	logHandler := slog.Default()

	// 尝试从 Locals 获取 traceId 和 spanId
	if traceId := c.Locals(TraceIdKey); traceId != nil {
		logHandler = logHandler.With(slog.String("traceId", ToString(traceId)))
	}
	if spanId := c.Locals(SpanIdKey); spanId != nil {
		logHandler = logHandler.With(slog.String("spanId", ToString(spanId)))
	}

	return logHandler
}

// propagateTraceIds 将 traceId 和 spanId 传播到响应头
func propagateTraceIds(c *fiber.Ctx) {
	if traceId := c.Locals(TraceIdKey); traceId != nil {
		c.Set(HeaderTrace, ToString(traceId))
	}
	if spanId := c.Locals(SpanIdKey); spanId != nil {
		c.Set(HeaderSpan, ToString(spanId))
	}
}
