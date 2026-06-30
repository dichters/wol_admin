package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mydev/mydev-api/enums"
	"github.com/mydev/mydev-api/response"
)

// DivideHandler 除法处理器
type DivideHandler struct{}

// NewDivideHandler 创建除法处理器
func NewDivideHandler() *DivideHandler {
	return &DivideHandler{}
}

// Handle 处理除法请求
// POST /srv/v1/divide
func (h *DivideHandler) Handle(c *fiber.Ctx) error {
	var req response.DivideRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusOK).JSON(response.Error(
			int(enums.InvalidParameter),
			"Invalid request body.",
		))
	}

	// 执行除法
	result := req.A / req.B

	return c.Status(fiber.StatusOK).JSON(response.SuccessWithData([]interface{}{result}))
}
