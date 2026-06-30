package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mydev/mydev-api/response"
)

// HealthHandler 健康检查处理器
type HealthHandler struct{}

// NewHealthHandler 创建健康检查处理器
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Handle 处理健康检查请求
// POST /srv/v1/hc
func (h *HealthHandler) Handle(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(response.Success())
}
